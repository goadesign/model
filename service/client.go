/*
Package service implements a client for the
[Structurizr service HTTP APIs](https://structurizr.com/).
*/
package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"goa.design/structurizr/expr"
)

// Host is the Structurizr API host.
const Host = "api.structurizr.com"

// UserAgent is the user agent used by this package.
const UserAgent = "structurizr-go/1.0"

// Response describes the API response returned by some endpoints.
type Response struct {
	// Success is true if the API call was successful, false otherwise.
	Success bool `json:"success"`
	// Message is a human readable response message.
	Message string `json:"message"`
	// Revision is hte internal revision number.
	Revision int `json:"revision"`
}

// Doer is an interface that encapsulate a HTTP client Do method.
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// Client holds the credentials needed to make the requests.
type Client struct {
	// Key is the API key.
	Key string
	// Secret is the API secret.
	Secret string
	// HTTP is the low level HTTP client.
	HTTP Doer
}

// NewClient instantiates a client using the default HTTP client.
func NewClient(key, secret string) *Client {
	return &Client{Key: key, Secret: secret, HTTP: http.DefaultClient}
}

// Get retrieves the given workspace.
func (c *Client) Get(id string) (*expr.Workspace, error) {
	u := &url.URL{Scheme: "https", Host: Host, Path: fmt.Sprintf("/workspace/%s", id)}
	req, _ := http.NewRequest("GET", u.String(), nil)
	c.sign(req, "", "")
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("service error: %s", string(body))
	}
	var workspace expr.Workspace
	if err := json.NewDecoder(resp.Body).Decode(&workspace); err != nil {
		return nil, err
	}
	return &workspace, nil
}

// Put stores the given workspace.
func (c *Client) Put(w *expr.Workspace) error {
	u := &url.URL{Scheme: "https", Host: Host, Path: fmt.Sprintf("/workspace/%s", w.ID)}
	body, _ := json.Marshal(w)
	req, _ := http.NewRequest("PUT", u.String(), bytes.NewReader(body))
	ct := "application/json; charset=UTF-8"
	c.sign(req, string(body), ct)
	req.Header.Set("Content-Type", ct)
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("service error: %s", string(body))
	}
	return nil
}

// Lock locks the given workspace.
func (c *Client) Lock(id string) (*Response, error) { return c.lockUnlock(id, true) }

// Unlock unlocks a previously locked workspace.
func (c *Client) Unlock(id string) (*Response, error) { return c.lockUnlock(id, false) }

// lockUnlock implements the Lock and Unlock calls.
func (c *Client) lockUnlock(id string, lock bool) (*Response, error) {
	u := &url.URL{Scheme: "https", Host: Host, Path: fmt.Sprintf("/workspace/%s/lock", id)}
	verb := "PUT"
	if !lock {
		verb = "DELETE"
	}
	req, _ := http.NewRequest(verb, u.String(), nil)
	c.sign(req, "", "")
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("service error: %s", string(body))
	}
	var res Response
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

// sign signs the requests as per https://structurizr.com/help/web-api
func (c *Client) sign(req *http.Request, content, ct string) {
	// 1. Compute digest
	var digest, nonce string
	{
		h := md5.New()
		h.Write([]byte(content))
		md5 := hex.EncodeToString(h.Sum(nil))
		nonce = strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
		digest = fmt.Sprintf("%s\n%s\n%s\n%s\n%s", req.Method, req.URL.Path, md5, ct, nonce)
	}

	// 2. Compute HMAC
	var hm []byte
	{
		h := hmac.New(sha256.New, []byte(c.Secret))
		h.Write([]byte(digest))
		hm = h.Sum(nil)
	}

	// 3. Write X-Authorization and Nonce headers
	req.Header.Set("X-Authorization", fmt.Sprintf("%s:%s", c.Key, base64.StdEncoding.EncodeToString(hm)))
	req.Header.Set("Nonce", nonce)

	// Finally set agent.
	req.Header.Set("User-Agent", UserAgent)
}
