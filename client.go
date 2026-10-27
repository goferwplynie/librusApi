package librus

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"
)

const (
	apiHost = "https://synergia.librus.pl/gateway/api/2.0/"
	// Entry point of the web login flow. Synergia redirects from here to the OAuth
	// authorization page with a state it later checks when the code comes back.
	portalLoginUrl = "https://synergia.librus.pl/loguj/portalRodzina"
	tokenInfoUrl   = apiHost + "Auth/TokenInfo/"
	maxAuthHops    = 10
)

// Client talks to the Librus Synergia API. It is safe for concurrent use.
//
// After a successful Login the client remembers the credentials and logs in
// again on its own when the session expires (Librus sessions last ~10-15 minutes).
type Client struct {
	mu         sync.Mutex
	httpClient *http.Client // nil until the first successful Login
	login      string
	password   string
}

// New creates a client. Call Login before fetching any resources.
func New() *Client {
	return &Client{}
}

// Login goes through the whole Librus auth flow and starts a new API session.
// The credentials are kept in memory so the client can log in again when the
// session expires.
func (c *Client) Login(ctx context.Context, login string, password string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	httpClient, err := newSession(ctx, login, password)
	if err != nil {
		return err
	}
	c.httpClient = httpClient
	c.login, c.password = login, password
	return nil
}

// Get sends a GET request to a raw API resource (e.g. "Subjects"), relative to
// https://synergia.librus.pl/gateway/api/2.0/. Use it for endpoints this
// library doesn't wrap yet. The caller must close the response body.
//
// Responses with status >= 400 are returned as *APIError. An expired session is
// renewed with the stored credentials and the request is retried once.
func (c *Client) Get(ctx context.Context, resource string) (*http.Response, error) {
	c.mu.Lock()
	httpClient := c.httpClient
	c.mu.Unlock()
	if httpClient == nil {
		return nil, ErrNotLoggedIn
	}

	resp, err := doGet(ctx, httpClient, apiHost+resource)
	var apiErr *APIError
	if !errors.As(err, &apiErr) || !apiErr.sessionExpired() {
		return resp, err
	}

	httpClient, err = c.relogin(ctx, httpClient)
	if err != nil {
		return nil, fmt.Errorf("librus: session expired and login failed: %w", err)
	}
	return doGet(ctx, httpClient, apiHost+resource)
}

// relogin starts a new session unless another goroutine already replaced stale.
func (c *Client) relogin(ctx context.Context, stale *http.Client) (*http.Client, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.httpClient != stale {
		return c.httpClient, nil
	}
	httpClient, err := newSession(ctx, c.login, c.password)
	if err != nil {
		return nil, err
	}
	c.httpClient = httpClient
	return httpClient, nil
}

// getJSON fetches resource and decodes its JSON body into v.
func (c *Client) getJSON(ctx context.Context, resource string, v any) error {
	resp, err := c.Get(ctx, resource)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return fmt.Errorf("librus: decoding %s response: %w", resource, err)
	}
	return nil
}

func doGet(ctx context.Context, httpClient *http.Client, target string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return nil, fmt.Errorf("librus: creating request: %w", err)
	}
	return do(httpClient, req)
}

// do sends req and turns responses with status >= 400 into *APIError.
func do(httpClient *http.Client, req *http.Request) (*http.Response, error) {
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("librus: request failed: %w", err)
	}
	if resp.StatusCode < 400 {
		return resp, nil
	}
	defer resp.Body.Close()

	apiErr := &APIError{StatusCode: resp.StatusCode}
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
	_ = json.Unmarshal(body, apiErr) // non-JSON error pages just leave Code and Message empty
	return nil, apiErr
}

// newSession logs in with a fresh cookie jar and returns an http.Client that
// carries the new session.
func newSession(ctx context.Context, login, password string) (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("librus: creating cookie jar: %w", err)
	}
	httpClient := &http.Client{Jar: jar}

	// redirects are followed manually so every hop stays on the same cookie jar
	auth := *httpClient
	auth.CheckRedirect = func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}

	authUrl, err := redirectTarget(ctx, &auth, portalLoginUrl)
	if err != nil {
		return nil, err
	}
	if authUrl == "" {
		return nil, errors.New("librus: login page didn't redirect to authorization")
	}

	// authorization page may redirect once more to the actual login form
	loginUrl, err := redirectTarget(ctx, &auth, authUrl)
	if err != nil {
		return nil, err
	}
	if loginUrl == "" {
		loginUrl = authUrl
	} else if _, err := redirectTarget(ctx, &auth, loginUrl); err != nil {
		return nil, err
	}

	data := url.Values{}
	data.Set("action", "login")
	data.Set("login", login)
	data.Set("pass", password)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, loginUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("librus: creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := do(&auth, req)
	var apiErr *APIError
	if errors.As(err, &apiErr) && apiErr.StatusCode < 500 {
		// wrong credentials come back as 403
		return nil, fmt.Errorf("%w: %w", ErrLoginFailed, err)
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var loginResp struct {
		Status string `json:"status"`
		GoTo   string `json:"goTo"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return nil, fmt.Errorf("librus: decoding login response: %w", err)
	}
	if loginResp.Status != "ok" || loginResp.GoTo == "" {
		return nil, fmt.Errorf("%w (status %q)", ErrLoginFailed, loginResp.Status)
	}

	next, err := resolveUrl(loginUrl, loginResp.GoTo)
	if err != nil {
		return nil, err
	}
	for hops := 0; next != ""; hops++ {
		if hops == maxAuthHops {
			return nil, errors.New("librus: authorization exceeded the redirect limit")
		}
		next, err = redirectTarget(ctx, &auth, next)
		if err != nil {
			return nil, err
		}
	}

	tokenResp, err := doGet(ctx, httpClient, tokenInfoUrl)
	if err != nil {
		return nil, fmt.Errorf("librus: login finished but api session wasn't created: %w", err)
	}
	tokenResp.Body.Close()

	return httpClient, nil
}

// redirectTarget sends GET to target and returns absolute Location of the response, or "" if it didn't redirect.
func redirectTarget(ctx context.Context, httpClient *http.Client, target string) (string, error) {
	resp, err := doGet(ctx, httpClient, target)
	if err != nil {
		return "", err
	}
	resp.Body.Close()

	location := resp.Header.Get("Location")
	if location == "" {
		return "", nil
	}
	return resolveUrl(target, location)
}

func resolveUrl(base, ref string) (string, error) {
	baseUrl, err := url.Parse(base)
	if err != nil {
		return "", fmt.Errorf("librus: invalid url %q: %w", base, err)
	}
	refUrl, err := baseUrl.Parse(ref)
	if err != nil {
		return "", fmt.Errorf("librus: invalid url %q: %w", ref, err)
	}
	return refUrl.String(), nil
}
