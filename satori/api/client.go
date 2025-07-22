package api

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

const satoriJwtFileName = "satori_jwt.txt"

func NewClient(host, userAgent, accountId, username, password *string, verifyTls bool, reuseJwt bool, jwtPath string) (*Client, error) {

	log.Printf("Creating a new client")

	if !verifyTls {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	jwtToken := ""

	if reuseJwt {
		jwtContent, err := loadJwtFromFile(jwtPath) //jwtPath is expected to be a directory or undefined
		if err != nil {
			// do nothing, just create a new JWT
			log.Printf("Failed to load file, creating a new token (%v)", err)
		} else {
			err := validateJwt(jwtContent)
			if err != nil {
				log.Printf("JWT validation error: %s , creating a new token", err)
				// JWT is not valid, create new one
			} else {
				// assign loaded JWT to client
				log.Printf("JWT validated and will be used")
				jwtToken = jwtContent
			}
		}
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

	if (jwtToken == "") && (username != nil) && (password != nil) {
		// form request body
		log.Printf("Authenticating and creating a new token")
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

		if reuseJwt {
			err := storeJwtToFile(jwtPath, ar.Token)
			if err != nil {
				log.Printf("JWT was not stored, %v", err)
			} else {
				log.Printf("JWT stored for future usage...")
			}
		}
	} else if jwtToken != "" {
		log.Printf("Reuse loaded JWT")
		c.Token = jwtToken
	}

	return &c, nil
}

func storeJwtToFile(predefinedPath string, token string) error {
	var filePath string

	if predefinedPath != "" {
		filePath = filepath.Join(predefinedPath, satoriJwtFileName)
	} else {
		tmpDir := os.TempDir()
		filePath = filepath.Join(tmpDir, satoriJwtFileName)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("failed to close file: %v", err)
		}
	}(file)

	_, err = file.WriteString(token)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	log.Printf("JWT written to: %s", filePath)
	return nil
}

func loadJwtFromFile(predefinedPath string) (string, error) {
	var filePath string

	if predefinedPath != "" {
		filePath = filepath.Join(predefinedPath, satoriJwtFileName)
	} else {
		tmpDir := os.TempDir()
		filePath = filepath.Join(tmpDir, satoriJwtFileName)
	}
	log.Printf("Loading JWT from: %s", filePath)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return string(data), nil
}

func validateJwt(tokenString string) error {
	// Define the claims structure
	claims := &jwt.RegisteredClaims{}

	// Parse the token without validating the signature
	_, _, err := jwt.NewParser().ParseUnverified(tokenString, claims)
	if err != nil {
		return fmt.Errorf("failed to parse token: %w", err)
	}

	// Validate expiration for 10 minutes from now
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now().Add(10*time.Minute)) {
		return fmt.Errorf("token is expired or will expire soon")
	}

	log.Printf("Token is valid for  usage")
	return nil
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

func (c *Client) getJsonById(path string, suffix string, id string, output interface{}) (error, int) {

	url := fmt.Sprintf("%s%s/%s", c.HostURL, path, id)
	if len(suffix) > 0 {
		url += "/" + suffix
	}

	req, err := http.NewRequest("GET", url, nil)
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

func (c *Client) getJsonForAccount(path string, search *string, output interface{}) (error, int) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", c.HostURL, path), nil)
	if err != nil {
		return err, 0
	}

	q := req.URL.Query()
	q.Add("accountId", c.AccountId)
	q.Add("page", "0")
	q.Add("pageSize", "200")
	if search != nil {
		q.Add("search", *search)
	}
	req.URL.RawQuery = q.Encode()

	body, err, responseStatus := c.doRequest(req)
	if err != nil {
		return err, responseStatus
	}

	err = json.Unmarshal(body, output)
	if err != nil {
		return err, responseStatus
	}

	return nil, responseStatus
}

func (c *Client) getJsonForAccountWithParams(path string, params *map[string]string, output interface{}) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", c.HostURL, path), nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	for k, v := range *params {
		q.Set(k, v)
	}
	q.Add("accountId", c.AccountId)
	q.Add("page", "0")
	q.Add("pageSize", "200")
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
		log.Printf("Recieved error: %s", err)
		return err
	}

	err = json.Unmarshal(body, output)
	if err != nil {
		log.Printf("Recieved unmarshal error: %s", err)
		return err
	}

	return nil
}

func (c *Client) putJson(path string, suffix string, id string, input interface{}, output interface{}) error {
	rb, err := json.Marshal(input)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s%s/%s", c.HostURL, path, id)
	if len(suffix) > 0 {
		url += "/" + suffix
	}

	req, err := http.NewRequest("PUT", url, strings.NewReader(string(rb)))
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

func (c *Client) putWithParams(path string, suffix string, id string, params *map[string]string, output interface{}) error {

	url := fmt.Sprintf("%s%s/%s", c.HostURL, path, id)
	if len(suffix) > 0 {
		url += "/" + suffix
	}

	req, err := http.NewRequest("PUT", url, nil)
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
