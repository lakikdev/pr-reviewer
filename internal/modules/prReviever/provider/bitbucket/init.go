package bitbucket

import (
	"fmt"
	"net/http"
	"net/url"
	"pr-reviewer/internal/modules/prReviever/provider/model"
)

type Client struct {
	apiBaseURL *url.URL
	auth       *auth
	client     *http.Client

	PullRequests PullRequestInterface
	DiffStat     DiffStatInterface
	Diff         DiffInterface
	File         FileInterface
	Comments     CommentsInterface
}

type auth struct {
	user string
	pass string
}

func NewClientBasicAuth(user, pass string) *Client {
	auth := &auth{
		user: user,
		pass: pass,
	}

	return buildClient(auth)
}

func buildClient(auth *auth) *Client {
	apiBaseURL, err := url.Parse("https://api.bitbucket.org/2.0/")
	if err != nil {
		panic(fmt.Sprintf("failed to parse API base URL: %v", err))
	}

	c := &Client{
		apiBaseURL: apiBaseURL,
		auth:       auth,
		client:     &http.Client{},
	}

	c.PullRequests = NewPullRequest(c)
	c.DiffStat = NewDiffStat(c)
	c.Diff = NewDiff(c)
	c.File = NewFile(c)
	c.Comments = NewComments(c)

	return c
}

type BitbucketPR struct {
	Owner         string `json:"owner"`
	RepoSlug      string `json:"repo_slug"`
	PullRequestID string `json:"pull_request_id"`
}

func (c *Client) GetPullRequest(owner, repoSlug, prID string) (*model.PullRequest, error) {
	pr, err := c.PullRequests.Get(&PullRequestOptions{
		ID:       prID,
		Owner:    owner,
		RepoSlug: repoSlug,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get pull request: %w", err)
	}

	return &model.PullRequest{
		ID:       prID,
		Owner:    owner,
		RepoSlug: repoSlug,

		DiffHref:     pr.Links.Diff.Href,
		DiffStatHref: pr.Links.DiffStat.Href,
	}, nil
}

func (c *Client) GetChangedFiles(href string) ([]model.FileChange, error) {
	diffStats, err := c.DiffStat.Get(&DiffStatOptions{
		Href: href,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get changed files: %w", err)
	}

	var files []model.FileChange
	for _, file := range diffStats.Values {
		files = append(files, model.FileChange{
			FilePath:    file.Old.Path,
			OldFileHref: file.Old.Links.Self.Href,
			NewFileHref: file.New.Links.Self.Href,
			Status:      file.Status,
		})
	}

	return files, nil
}

func (c *Client) GetFileDiff(filePath string, diffHref string) (string, error) {
	diff, err := c.Diff.Get(&DiffOptions{
		Href: diffHref,
		Path: filePath,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get file diff: %w", err)
	}
	return diff.Diff, nil
}

func (c *Client) LoadFile(href string) (string, error) {
	fileContent, err := c.File.Get(&FileOptions{
		Href: href,
	})
	if err != nil {
		return "", fmt.Errorf("failed to load file: %w", err)
	}
	return fileContent, nil
}

func (c *Client) AddComment(owner, repoSlug, prID string, filePath string, line int, comment string) error {
	if err := c.Comments.Add(&CommentsOptions{
		Owner:         owner,
		RepoSlug:      repoSlug,
		PullRequestID: prID,
		FilePath:      filePath,
		LineNumber:    line,
		Comment:       comment,
	}); err != nil {
		return fmt.Errorf("failed to add comment: %w", err)
	}
	return nil
}
