package provider

import (
	"fmt"
	"pr-reviewer/internal/config"
	"pr-reviewer/internal/modules/prReviever/provider/bitbucket"
	"pr-reviewer/internal/modules/prReviever/provider/model"
)

type Interface interface {
	GetPullRequest(owner, repoSlug, prID string) (*model.PullRequest, error)
	GetChangedFiles(href string) ([]model.FileChange, error)
	GetFileDiff(filePath string, diffHref string) (string, error)
	LoadFile(href string) (string, error)
	AddComment(owner, repoSlug, prID string, filePath string, line int, comment string) error
}

func GetProvider(name model.ProviderName) (Interface, error) {
	switch name {
	case model.Bitbucket:
		return bitbucket.NewClientBasicAuth(*config.BitbucketUsername, *config.BitbucketPassword), nil
	default:
		return nil, fmt.Errorf("provider %s not supported", name)
	}
}
