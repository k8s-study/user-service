package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
)

// Client is RESTAPI client with tracing header injection
type Client struct {
	BaseURL         *url.URL
	httpClient      *http.Client
	originalRequest *http.Request
}

func kongHostURL() (url string) {
	url = os.Getenv("KONG_HOST")
	if url == "" {
		url = "http://kong-ingress-controller.kong:8001"
	}

	return
}

// NewClient returns a new REST API client
func NewClient(req *http.Request) *Client {
	if u, err := url.Parse(kongHostURL()); err == nil {
		c := &Client{BaseURL: u, httpClient: &http.Client{}, originalRequest: req}
		return c
	}

	return nil
}

// NewRequest is create request (not call)
func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	// set tracing header
	incomingHeaders := []string{
		"x-request-id",
		"x-b3-traceid",
		"x-b3-spanid",
		"x-b3-parentspanid",
		"x-b3-sampled",
		"x-b3-flags",
		"x-ot-span-context",
	}

	for _, h := range incomingHeaders {
		req.Header.Set(h, c.originalRequest.Header.Get(h))
	}

	return req, nil
}

// Do is call request
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}
