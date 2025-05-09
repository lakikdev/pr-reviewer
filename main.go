package main

import (
	prReviewer "pr-reviewer/internal/modules/prReviever"
	aiModel "pr-reviewer/internal/modules/prReviever/ai/model"
	providerModel "pr-reviewer/internal/modules/prReviever/provider/model"
)

//test comment

const (
	BITBUCKET_REPO_OWNER = "lakik"
	BITBUCKET_REPO_SLUG  = "pr-reviewer"
	BITBUCKET_PR_ID      = "1"
)

func main() {
	reviewer, err := prReviewer.NewReviewer(&prReviewer.ReviewerOptions{
		Provider: providerModel.Bitbucket,
		AIClient: aiModel.Ollama,
	})
	if err != nil {
		panic(err)
	}

	// Analyze the pull request
	commentsAdded, err := reviewer.Analyze(BITBUCKET_REPO_OWNER, BITBUCKET_REPO_SLUG, BITBUCKET_PR_ID)
	if err != nil {
		panic(err)
	}

	// Print the number of comments added
	println("Number of comments added:", commentsAdded)
}
