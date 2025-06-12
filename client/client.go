package client

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

type Client struct {
	Host       string
	Logged_in  bool
	HttpClient *http.Client
}

func New() *Client {
	jar, _ := cookiejar.New(nil)
	return &Client{
		Host:      "https://synergia.librus.pl/gateway/api/2.0/",
		Logged_in: false,
		HttpClient: &http.Client{
			Jar: jar,
		},
	}
}

func (c *Client) Login(login string, password string) error {

	testauthUrl := "https://api.librus.pl/OAuth/Authorization?client_id=46&response_type=code&scope=mydata"
	grantUrl := "https://api.librus.pl/OAuth/Authorization/Grant?client_id=46"
	tokenInfo := "https://synergia.librus.pl/gateway/api/2.0/Auth/TokenInfo/"

	req, err := http.NewRequest(http.MethodGet, testauthUrl, nil)

	if err != nil {
		return fmt.Errorf("cant create request:\n%w", err)
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed:\n%w", err)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("not nice response code :c :%v", resp.StatusCode)
	}

	loginEndpoint := "https://api.librus.pl/OAuth/Authorization?client_id=46"
	data := url.Values{}
	data.Set("action", "login")
	data.Set("login", login)
	data.Set("pass", password)

	req, err = http.NewRequest(http.MethodPost, loginEndpoint, strings.NewReader(data.Encode()))

	if err != nil {
		return fmt.Errorf("cant create request:\n%w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err = c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed:\n%w", err)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("not nice response code :c :%v", resp.StatusCode)
	}

	req, _ = http.NewRequest(http.MethodGet, grantUrl, nil)

	c.HttpClient.Do(req)

	req, _ = http.NewRequest(http.MethodGet, tokenInfo, nil)
	c.HttpClient.Do(req)

	return nil
}

func (c *Client) Get(resource string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, c.Host+resource, nil)

	if err != nil {
		return nil, err
	}

	res, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, err
	}

	return res, nil
}
