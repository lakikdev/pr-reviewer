package ai

import (
	"pr-reviewer/internal/modules/prReviever/ai/model"
	"pr-reviewer/internal/modules/prReviever/ai/ollama"
)

type ClientInterface interface {
	Send(prompt string) (string, error)
	ParseResponse(response string) model.StructResponse
}

func GetAIClient(name model.AIName) (ClientInterface, error) {
	switch name {
	case model.Ollama:
		return ollama.NewClient(), nil
	default:
		return nil, nil
	}
}
