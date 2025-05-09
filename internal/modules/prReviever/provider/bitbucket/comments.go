package bitbucket

import (
	"encoding/json"
	"fmt"
)

type Comments struct {
	c *Client
}

type CommentsInterface interface {
	Add(opts *CommentsOptions) error
	Get(opts *CommentsOptions) ([]CommentsResponse, error)
}

func NewComments(c *Client) CommentsInterface {
	return &Comments{
		c: c,
	}
}

type CommentsOptions struct {
	Owner         string
	RepoSlug      string
	PullRequestID string
	FilePath      string
	LineNumber    int
	Comment       string
}

func (c *Comments) Add(opts *CommentsOptions) error {
	payload := map[string]interface{}{
		"content": map[string]string{
			"raw": fmt.Sprintf("<PR-Reviewer Bot> - %s", opts.Comment),
		},
		"inline": map[string]interface{}{
			"path": opts.FilePath,
			"to":   opts.LineNumber,
		},
	}

	req, err := c.c.newRequest("POST", "/repositories/"+opts.Owner+"/"+opts.RepoSlug+"/pullrequests/"+opts.PullRequestID+"/comments", payload)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	response, err := c.c.doRequest(req)
	if err != nil {
		return fmt.Errorf("failed to add comment: %w", err)
	}

	if response.StatusCode != 201 {
		return fmt.Errorf("unexpected response status: %d", response.StatusCode)
	}

	return nil
}

type CommentsResponse struct {
	Content struct {
		Raw string `json:"raw"`
	} `json:"content"`
	Inline struct {
		Path string `json:"path"`
		To   int    `json:"to"`
	} `json:"inline"`
}

func (c *Comments) Get(opts *CommentsOptions) ([]CommentsResponse, error) {
	req, err := c.c.newRequest("GET", "/repositories/"+opts.Owner+"/"+opts.RepoSlug+"/pullrequests/"+opts.PullRequestID+"/comments", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		//TODO Handle non-200 status codes
		return nil, fmt.Errorf("failed to get comments: %s", resp.Status)
	}

	var comments []CommentsResponse
	if err := json.NewDecoder(resp.Body).Decode(&comments); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return comments, nil
}
