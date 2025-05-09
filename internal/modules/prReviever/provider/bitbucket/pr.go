package bitbucket

import (
	"encoding/json"
	"fmt"
)

type PullRequest struct {
	c *Client
}

type PullRequestInterface interface {
	Get(opts *PullRequestOptions) (*PullRequestResponse, error)
}

type PullRequestOptions struct {
	ID       string
	Owner    string
	RepoSlug string
}

type PullRequestResponse struct {
	Links struct {
		Diff struct {
			Href string `json:"href"`
		} `json:"diff"`
		DiffStat struct {
			Href string `json:"href"`
		} `json:"diffstat"`
	} `json:"links"`
}

func NewPullRequest(c *Client) PullRequestInterface {
	return &PullRequest{
		c: c,
	}
}

func (p *PullRequest) Get(opts *PullRequestOptions) (*PullRequestResponse, error) {
	req, err := p.c.newRequest("GET", "/repositories/"+opts.Owner+"/"+opts.RepoSlug+"/pullrequests/"+opts.ID, nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		//TODO Handle non-200 status codes
		return nil, fmt.Errorf("failed to get pull request: %s", resp.Status)
	}

	
	var pr PullRequestResponse
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check if the response contains the expected fields
	if pr.Links.Diff.Href == "" || pr.Links.DiffStat.Href == "" {
		return nil, fmt.Errorf("unexpected response format: %v", pr)
	}

	return &pr, nil
}
