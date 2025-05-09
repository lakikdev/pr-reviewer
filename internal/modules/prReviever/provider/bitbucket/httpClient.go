package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	if c.apiBaseURL == nil {
		return nil, fmt.Errorf("apiBaseURL is nil")
	}

	//if path is full url, use it as is
	var requestURL *url.URL
	if strings.HasPrefix(path, "http") {
		var err error
		requestURL, err = url.Parse(path)
		if err != nil {
			return nil, fmt.Errorf("parsing URL: %w", err)
		}
	} else {
		path = strings.TrimPrefix(path, "/")
		requestURL = c.apiBaseURL.ResolveReference(&url.URL{Path: path})
	}

	// Debugging line
	fmt.Printf("API Base URL: %s\n", c.apiBaseURL.String()) // Debugging line
	fmt.Printf("path: %s\n", path)                          // Debugging line
	fmt.Printf("Request URL: %s\n", requestURL.String())    // Debugging line

	req, err := http.NewRequest(method, requestURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	if body != nil {
		err = c.encodeBody(req, body)
		if err != nil {
			return nil, fmt.Errorf("encoding body: %w", err)
		}
	}

	return req, nil
}

func (c *Client) encodeBody(req *http.Request, body interface{}) error {
	if body == nil {
		return nil
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshalling body: %w", err)
	}

	req.Body = io.NopCloser(bytes.NewBuffer(jsonBody))
	req.ContentLength = int64(len(jsonBody))

	return nil
}

func (c *Client) doRequest(req *http.Request) (*http.Response, error) {
	if req == nil {
		return nil, fmt.Errorf("request is nil")
	}

	if c.auth.pass != "" && c.auth.user != "" {
		req.SetBasicAuth(c.auth.user, c.auth.pass)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}

	return resp, nil
}
