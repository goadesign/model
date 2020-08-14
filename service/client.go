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
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/user"
	"strconv"
	"time"

	goahttp "goa.design/goa/v3/http"
	"goa.design/model/expr"
)

var (
	// Host is the Structurizr API host (var for testing).
	Host = "api.structurizr.com"

	// Scheme is the HTTP scheme used to make requests to the Structurizr service.
	Scheme = "https"
)

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
	u := &url.URL{Scheme: Scheme, Host: Host, Path: fmt.Sprintf("/workspace/%s", id)}
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
func (c *Client) Put(id string, w *expr.Workspace) error {
	u := &url.URL{Scheme: Scheme, Host: Host, Path: fmt.Sprintf("/workspace/%s", id)}
	body, _ := json.Marshal(w)
	req, _ := http.NewRequest("PUT", u.String(), bytes.NewReader(body))
	ct := "application/json; charset=UTF-8"
	c.sign(req, string(body), ct)
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
func (c *Client) Lock(id string) error { return c.lockUnlock(id, true) }

// Unlock unlocks a previously locked workspace.
func (c *Client) Unlock(id string) error { return c.lockUnlock(id, false) }

// EnableDebug causes the client to print debug information to Stderr.
func (c *Client) EnableDebug() {
	c.HTTP = goahttp.NewDebugDoer(c.HTTP)
}

// lockUnlock implements the Lock and Unlock calls.
func (c *Client) lockUnlock(id string, lock bool) error {
	u := &url.URL{Scheme: Scheme, Host: Host, Path: fmt.Sprintf("/workspace/%s/lock", id)}
	name := "anon"
	if usr, err := user.Current(); err == nil {
		name = usr.Name
		if name == "" {
			name = usr.Username
		}
	}
	// the order matters :(
	u.RawQuery = "user=" + url.QueryEscape(name) + "&agent=" + url.QueryEscape(UserAgent)

	verb := "PUT"
	if !lock {
		verb = "DELETE"
	}
	req, _ := http.NewRequest(verb, u.String(), nil)
	c.sign(req, "", "")
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var res Response
		json.NewDecoder(resp.Body).Decode(&res) // ignore error, just trying
		err = fmt.Errorf("service error: %s", resp.Status)
		if res.Message != "" {
			err = errors.New(res.Message)
		}
		return err
	}

	return nil
}

// sign signs the requests as per https://structurizr.com/help/web-api
func (c *Client) sign(req *http.Request, content, ct string) {
	// 1. Compute digest
	var digest, nonce string
	var md5s string
	{
		h := md5.New()
		h.Write([]byte(content))
		md5b := h.Sum(nil)
		md5s = hex.EncodeToString(md5b)
		nonce = strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
		digest = fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n", req.Method, req.URL.RequestURI(), md5s, ct, nonce)
	}

	// 2. Compute HMAC
	var hm []byte
	{
		h := hmac.New(sha256.New, []byte(c.Secret))
		h.Write([]byte(digest))
		hm = h.Sum(nil)
	}

	// 3. Write headers
	auth := base64.StdEncoding.EncodeToString([]byte(hex.EncodeToString(hm)))
	req.Header.Set("X-Authorization", fmt.Sprintf("%s:%s", c.Key, auth))
	req.Header.Set("Nonce", nonce)
	if req.Method == http.MethodPut {
		req.Header.Set("Content-Md5", base64.StdEncoding.EncodeToString([]byte(md5s)))
		req.Header.Set("Content-Type", ct)
	}
	req.Header.Set("User-Agent", UserAgent)
}
