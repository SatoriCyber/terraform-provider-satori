package api

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// HostURL - Default Satori URL
const HostURL string = "https://app.satoricyber.com"

type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
	AccountId  string
	UserAgent  string
}

type AuthStruct struct {
	Username string `json:"serviceAccountId"`
	Password string `json:"serviceAccountKey"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func NewClient(host, userAgent, accountId, username, password *string, verifyTls bool) (*Client, error) {

	if !verifyTls {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		HostURL:    HostURL,
		AccountId:  *accountId,
		UserAgent:  *userAgent,
	}

	if host != nil {
		c.HostURL = *host
	}

	if (username != nil) && (password != nil) {
		// form request body
		rb, err := json.Marshal(AuthStruct{
			Username: *username,
			Password: *password,
		})
		if err != nil {
			return nil, err
		}

		// authenticate
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/authentication/token", c.HostURL), strings.NewReader(string(rb)))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")

		body, err, _ := c.doRequest(req)
		if err != nil {
			return nil, err
		}

		// parse response body
		ar := AuthResponse{}
		err = json.Unmarshal(body, &ar)
		if err != nil {
			return nil, err
		}

		c.Token = ar.Token
	}

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error, int) {
	if len(c.Token) > 0 {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	req.Header.Set("User-Agent", c.UserAgent)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err, 0
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err, 0
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body), res.StatusCode
	}

	return body, err, res.StatusCode
}

func (c *Client) getJsonById(path string, id string, output interface{}) (error, int) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s/%s", c.HostURL, path, id), nil)
	if err != nil {
		return err, 0
	}

	body, err, statusCode := c.doRequest(req)
	if err != nil {
		return err, statusCode
	}

	err = json.Unmarshal(body, output)
	if err != nil {
		return err, statusCode
	}

	return nil, statusCode
}

func (c *Client) getJsonForAccount(path string, search *string, output interface{}) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", c.HostURL, path), nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("accountId", c.AccountId)
	q.Add("page", "0")
	q.Add("pageSize", "200")
	if search != nil {
		q.Add("search", *search)
	}
	req.URL.RawQuery = q.Encode()

	body, err, _ := c.doRequest(req)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, output)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) postJsonForAccount(path string, input interface{}, output interface{}) error {
	params := make(map[string]string, 1)
	params["accountId"] = c.AccountId
	return c.postJsonWithParams(path, &params, input, output)
}

func (c *Client) postJsonWithParams(path string, params *map[string]string, input interface{}, output interface{}) error {
	rb, err := json.Marshal(input)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", c.HostURL, path), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()
	for k, v := range *params {
		q.Set(k, v)
	}
	req.URL.RawQuery = q.Encode()

	body, err, _ := c.doRequest(req)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, output)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) putJson(path string, id string, input interface{}, output interface{}) error {
	rb, err := json.Marshal(input)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s%s/%s", c.HostURL, path, id), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	body, err, _ := c.doRequest(req)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, output)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) putWithParams(path string, id string, params *map[string]string, output interface{}) error {
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s%s/%s", c.HostURL, path, id), nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	for k, v := range *params {
		q.Set(k, v)
	}
	req.URL.RawQuery = q.Encode()

	body, err, _ := c.doRequest(req)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, output)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) putJsonWithParams(path string, id string, params *map[string]string, input interface{}, output interface{}) error {
	rb, err := json.Marshal(input)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s%s/%s", c.HostURL, path, id), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	q := req.URL.Query()
	for k, v := range *params {
		q.Set(k, v)
	}
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Content-Type", "application/json")

	body, err, _ := c.doRequest(req)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, output)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) delete(path string, id string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s%s/%s", c.HostURL, path, id), nil)
	if err != nil {
		return err
	}

	_, err, _ = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
