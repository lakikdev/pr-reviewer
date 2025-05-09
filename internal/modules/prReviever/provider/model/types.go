package model

type ProviderName string

const (
	Bitbucket ProviderName = "bitbucket"
	Github    ProviderName = "github"
	Gitlab    ProviderName = "gitlab"
)

type PullRequest struct {
	ID       string
	Owner    string
	RepoSlug string

	DiffHref     string
	DiffStatHref string
}

type FileChange struct {
	FilePath    string
	OldFileHref string
	NewFileHref string
	Status      string
}
