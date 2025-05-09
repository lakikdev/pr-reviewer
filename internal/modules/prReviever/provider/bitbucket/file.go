package bitbucket

import (
	"fmt"
	"io"
)

type File struct {
	c *Client
}

type FileInterface interface {
	Get(path *FileOptions) (string, error)
}

func NewFile(c *Client) FileInterface {
	return &File{
		c: c,
	}
}

type FileOptions struct {
	Href string
}

func (f *File) Get(opts *FileOptions) (string, error) {
	req, err := f.c.newRequest("GET", opts.Href, nil)
	if err != nil {
		return "", err
	}

	resp, err := f.c.doRequest(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to get file: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}
