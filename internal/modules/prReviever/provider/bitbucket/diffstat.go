package bitbucket

import (
	"encoding/json"
	"fmt"
)

type DiffStat struct {
	c *Client
}

type DiffStatInterface interface {
	Get(opts *DiffStatOptions) (*DiffStatResponse, error)
}

type DiffStatOptions struct {
	Href string
}

type DiffStatResponse struct {
	Values []struct {
		Status string `json:"status"`
		Old    struct {
			Path  string `json:"path"`
			Links struct {
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
			} `json:"links"`
		} `json:"old"`
		New struct {
			Path  string `json:"path"`
			Links struct {
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
			} `json:"links"`
		} `json:"new"`
	} `json:"values"`
	PageLen int    `json:"pagelen"`
	Next    string `json:"next"`
	Size    int    `json:"size"`
	Page    int    `json:"page"`
}

func NewDiffStat(c *Client) DiffStatInterface {
	return &DiffStat{
		c: c,
	}
}

func (d *DiffStat) Get(opts *DiffStatOptions) (*DiffStatResponse, error) {
	req, err := d.c.newRequest("GET", opts.Href, nil)
	if err != nil {
		return nil, err
	}

	resp, err := d.c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		//TODO Handle non-200 status codes
		return nil, fmt.Errorf("failed to get diffstat: %s", resp.Status)
	}

	var diffStat DiffStatResponse
	if err := json.NewDecoder(resp.Body).Decode(&diffStat); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &diffStat, nil
}
