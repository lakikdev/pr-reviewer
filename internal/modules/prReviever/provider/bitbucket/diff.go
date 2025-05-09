package bitbucket

import (
	"fmt"
	"io"
)

type Diff struct {
	c *Client
}

type DiffInterface interface {
	Get(opts *DiffOptions) (*DiffResponse, error)
}

func NewDiff(c *Client) DiffInterface {
	return &Diff{
		c: c,
	}
}

type DiffOptions struct {
	Href string
	Path string
}

type DiffResponse struct {
	Diff string `json:"diff"`
}

func (d *Diff) Get(opts *DiffOptions) (*DiffResponse, error) {
	req, err := d.c.newRequest("GET", opts.Href, nil)
	if err != nil {
		return nil, err
	}

	if opts.Path != "" {
		// Add the path to the request URL
		q := req.URL.Query()
		q.Add("path", opts.Path)
		req.URL.RawQuery = q.Encode()
	}

	resp, err := d.c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get diff: %s", resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	// Print the response body for debugging
	//fmt.Printf("Response Body: %s\n", string(body))

	// var diff DiffResponse
	// if err := json.NewDecoder(resp.Body).Decode(&diff); err != nil {
	// 	return nil, fmt.Errorf("failed to decode response: %w", err)
	// }

	return &DiffResponse{
		Diff: string(body),
	}, nil
}
