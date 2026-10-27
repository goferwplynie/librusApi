package librus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

const (
	// Entry point of the web login flow. Synergia redirects from here to the OAuth
	// authorization page with a state it later checks when the code comes back.
	portalLoginUrl = "https://synergia.librus.pl/loguj/portalRodzina"
	tokenInfoUrl   = "https://synergia.librus.pl/gateway/api/2.0/Auth/TokenInfo/"
	maxAuthHops    = 10
)

type Client struct {
	host       string
	loggedIn   bool
	httpClient *http.Client
}

// New creates a client. Call Login before fetching any resources.
func New() *Client {
	jar, _ := cookiejar.New(nil)
	return &Client{
		host:     "https://synergia.librus.pl/gateway/api/2.0/",
		loggedIn: false,
		httpClient: &http.Client{
			Jar: jar,
		},
	}
}

func (c *Client) request(method, target string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, target, body)

	if err != nil {
		return nil, fmt.Errorf("cant create request:\n%w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed:\n%w", err)
	}
	if resp.StatusCode >= 400 {
		return resp, fmt.Errorf("not nice response code :c :%v", resp.StatusCode)
	}

	return resp, nil
}

// login to librus api. passes whole auth flow. doesn't save login info so after session expired error you have to login again
func (c *Client) Login(login string, password string) error {
	// redirects are followed manually so every hop stays on the same cookie jar
	auth := *c.httpClient
	auth.CheckRedirect = func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}

	authUrl, err := redirectTarget(&auth, portalLoginUrl)
	if err != nil {
		return err
	}
	if authUrl == "" {
		return fmt.Errorf("librus login page didn't redirect to authorization")
	}

	// authorization page may redirect once more to the actual login form
	loginUrl, err := redirectTarget(&auth, authUrl)
	if err != nil {
		return err
	}
	if loginUrl == "" {
		loginUrl = authUrl
	} else if _, err := redirectTarget(&auth, loginUrl); err != nil {
		return err
	}

	data := url.Values{}
	data.Set("action", "login")
	data.Set("login", login)
	data.Set("pass", password)

	resp, err := auth.PostForm(loginUrl, data)
	if err != nil {
		return fmt.Errorf("request failed:\n%w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("not nice response code :c :%v", resp.StatusCode)
	}

	var loginResp struct {
		Status string `json:"status"`
		GoTo   string `json:"goTo"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return fmt.Errorf("failed to unmarshal login response:\n%w", err)
	}
	if loginResp.Status != "ok" || loginResp.GoTo == "" {
		return fmt.Errorf("login failed (status %q), check login and password", loginResp.Status)
	}

	next, err := resolveUrl(loginUrl, loginResp.GoTo)
	if err != nil {
		return err
	}
	for hops := 0; next != ""; hops++ {
		if hops == maxAuthHops {
			return fmt.Errorf("librus authorization exceeded the redirect limit")
		}
		next, err = redirectTarget(&auth, next)
		if err != nil {
			return err
		}
	}

	tokenResp, err := c.request(http.MethodGet, tokenInfoUrl, nil)
	if tokenResp != nil {
		tokenResp.Body.Close()
	}
	if err != nil {
		return fmt.Errorf("login finished but api session wasn't created:\n%w", err)
	}

	c.loggedIn = true
	return nil
}

// redirectTarget sends GET to target and returns absolute Location of the response, or "" if it didn't redirect.
func redirectTarget(client *http.Client, target string) (string, error) {
	resp, err := client.Get(target)
	if err != nil {
		return "", fmt.Errorf("request failed:\n%w", err)
	}
	resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("not nice response code :c :%v", resp.StatusCode)
	}

	location := resp.Header.Get("Location")
	if location == "" {
		return "", nil
	}
	return resolveUrl(target, location)
}

func resolveUrl(base, ref string) (string, error) {
	baseUrl, err := url.Parse(base)
	if err != nil {
		return "", fmt.Errorf("invalid url %q:\n%w", base, err)
	}
	refUrl, err := baseUrl.Parse(ref)
	if err != nil {
		return "", fmt.Errorf("invalid url %q:\n%w", ref, err)
	}
	return refUrl.String(), nil
}

// Get sends a GET request to a raw API resource (e.g. "Subjects"), relative to
// https://synergia.librus.pl/gateway/api/2.0/. Use it for endpoints this
// library doesn't wrap yet. The caller must close the response body.
func (c *Client) Get(resource string) (*http.Response, error) {
	resp, err := c.request(http.MethodGet, c.host+resource, nil)
	if resp != nil && resp.StatusCode == http.StatusUnauthorized {
		resp.Body.Close()
		c.loggedIn = false
		return nil, SessionExpiredError{}
	}
	return resp, err
}
