package librusApi

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

const (
	testauthUrl   = "https://api.librus.pl/OAuth/Authorization?client_id=46&response_type=code&scope=mydata"
	grantUrl      = "https://api.librus.pl/OAuth/Authorization/Grant?client_id=46"
	tokenInfo     = "https://synergia.librus.pl/gateway/api/2.0/Auth/TokenInfo/"
	loginEndpoint = "https://api.librus.pl/OAuth/Authorization?client_id=46"
)

type Client struct {
	host       string
	logged_in  bool
	httpClient *http.Client
}

// create new client
func New() *Client {
	jar, _ := cookiejar.New(nil)
	return &Client{
		host:      "https://synergia.librus.pl/gateway/api/2.0/",
		logged_in: false,
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

	_, err := c.request(http.MethodGet, testauthUrl, nil)
	if err != nil {
		return err
	}

	data := url.Values{}
	data.Set("action", "login")
	data.Set("login", login)
	data.Set("pass", password)

	req, err := http.NewRequest(http.MethodPost, loginEndpoint, strings.NewReader(data.Encode()))

	if err != nil {
		return fmt.Errorf("cant create request:\n%w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed:\n%w", err)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("not nice response code :c :%v", resp.StatusCode)
	}

	_, err = c.request(http.MethodGet, grantUrl, nil)
	if err != nil {
		return err
	}

	_, err = c.request(http.MethodGet, tokenInfo, nil)
	if err != nil && resp == nil {
		return err
	}

	return nil
}

// Helper function for getting resource from api.
func (c *Client) Get(resource string, body io.Reader) (*http.Response, error) {
	resp, err := c.request(http.MethodGet, c.host+resource, body)
	if resp.StatusCode == 401 {
		c.logged_in = false
		return nil, SessionExpiredError{}
	}
	return resp, err
}

func WithReauth[T any](
	method func() (T, error),
	loginFunc func() error,
	retries int,
) (T, error) {
	var zeroVal T
	for range retries + 1 {
		resp, err := method()
		if err == nil {
			return resp, nil
		}

		if IsSessionExpired(err) {
			err := loginFunc()
			if err != nil {
				return resp, fmt.Errorf("error in login func: %v", err)
			}
			continue
		}
		return zeroVal, err
	}
	return zeroVal, fmt.Errorf("reauth failed")
}
