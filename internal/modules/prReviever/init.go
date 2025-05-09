package prReviewer

import (
	"fmt"
	"pr-reviewer/internal/modules/prReviever/ai"
	aiModel "pr-reviewer/internal/modules/prReviever/ai/model"
	"pr-reviewer/internal/modules/prReviever/provider"
	providerModel "pr-reviewer/internal/modules/prReviever/provider/model"
	"strings"
)

type ReviewerOptions struct {
	Provider providerModel.ProviderName
	AIClient aiModel.AIName
}

type Reviewer struct {
	PullRequest *providerModel.PullRequest

	provider provider.Interface
	aiClient ai.ClientInterface
}

func NewReviewer(options *ReviewerOptions) (*Reviewer, error) {
	providerClient, err := provider.GetProvider(options.Provider)
	if err != nil {
		return nil, err
	}

	aiClient, err := ai.GetAIClient(options.AIClient)
	if err != nil {
		return nil, err
	}

	return &Reviewer{
		provider: providerClient,
		aiClient: aiClient,
	}, nil
}

func (r *Reviewer) Analyze(owner, repoSlug, pullRequestID string) (int, error) {
	fmt.Printf("Analyzing pull request %s/%s/%s\n", owner, repoSlug, pullRequestID)

	// Get the pull request
	pr, err := r.provider.GetPullRequest(owner, repoSlug, pullRequestID)
	if err != nil {
		return 0, err
	}

	r.PullRequest = pr

	// Get the diff stat
	changedFiles, err := r.provider.GetChangedFiles(pr.DiffStatHref)
	if err != nil {
		return 0, err
	}

	commentsAdded := 0
	for _, file := range changedFiles {
		// Get the diff
		diff, err := r.provider.GetFileDiff(file.FilePath, pr.DiffHref)
		if err != nil {
			return 0, err
		}

		// Get the file content
		fileContent, err := r.provider.LoadFile(file.NewFileHref)
		if err != nil {
			return 0, err
		}

		//prepare the prompt
		prompt := preparePrompt(diff, fileContent)

		// Send the diff and file content to the AI client

		response, err := r.aiClient.Send(prompt)
		if err != nil {
			return 0, err
		}
		// Parse the AI response
		parsedResponse := r.aiClient.ParseResponse(response)

		for _, comment := range parsedResponse.Comments {
			fmt.Printf("Line: %d Comment: %s\n", comment.LineNumber, comment.Text)

			//TODO to make sure we don't post the same comment multiple times
			//TODO add filtration to comments to make sure we don't post unnecessary comments

			// Post the comment to Bitbucket
			if err := r.provider.AddComment(owner, repoSlug, pullRequestID, file.FilePath, comment.LineNumber, comment.Text); err != nil {
				fmt.Printf("Error posting comment: %s\n", err)
				continue
			}

			commentsAdded++
		}

	}
	return commentsAdded, nil
}

func preparePrompt(diff, code string) string {
	var fullBody strings.Builder

	//TODO map diff to the code and mark what lines are changed
	fullBody.WriteString(diff)

	var codeWithLines strings.Builder
	for i, line := range strings.Split(code, "\n") {
		// Add line numbers to the code
		codeWithLines.WriteString(fmt.Sprintf("Line %d: %s\n", i+1, line))
	}

	fullBody.WriteString(codeWithLines.String())
	fullBody.WriteString("\n\nPlease review it.")
	fmt.Printf("Full Body: %s\n", fullBody.String()) // Debugging line

	return fullBody.String()

}
