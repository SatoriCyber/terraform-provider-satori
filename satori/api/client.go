package api

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/facette/logger"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// HostURL - Default Satori URL
const HostURL string = "https://app.satoricyber.com"

// Client -
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
	AccountId  string
	Logger     *logger.Logger
}

// AuthStruct -
type AuthStruct struct {
	Username string `json:"serviceAccountId"`
	Password string `json:"serviceAccountKey"`
}

// AuthResponse -
type AuthResponse struct {
	Token string `json:"token"`
}

// NewClient -
func NewClient(host, accountId, username, password *string, verifyTls bool) (*Client, error) {

	newLogger, err := logger.NewLogger(
		logger.FileConfig{
			Level: "debug",
			//Path:  "/Users/alex/workspace/terraform-provider-satoricyber/examples/log.log",
		},
	)
	if err != nil {
		log.Fatalf("failed to initialize logger: %s", err)
	}

	newLogger.Info("test test test")
	if !verifyTls {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		HostURL:    HostURL,
		AccountId:  *accountId,
		Logger:     newLogger,
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

		body, err := c.doRequest(req)
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

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	if len(c.Token) > 0 {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}

func (c *Client) getJson(path string, id string, output interface{}) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s/%s", c.HostURL, path, id), nil)
	if err != nil {
		return err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, output)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) getJsonForAccount(path string, output interface{}) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", c.HostURL, path), nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("accountId", c.AccountId)
	q.Add("page", "0")
	q.Add("pageSize", "200")
	req.URL.RawQuery = q.Encode()

	body, err := c.doRequest(req)
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
	q.Add("accountId", c.AccountId)
	req.URL.RawQuery = q.Encode()

	body, err := c.doRequest(req)
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

	body, err := c.doRequest(req)
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

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
