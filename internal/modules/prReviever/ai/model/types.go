package model

type AIName string

const (
	OpenAI AIName = "openai"
	Ollama AIName = "ollama"
	Google AIName = "google"
)

// SupportedProviders is a map of supported AI providers
var SupportedProviders = map[AIName]struct{}{
	Ollama: {},
}

type StructResponse struct {
	Comments []struct {
		LineNumber int    `json:"line_number"`
		Text       string `json:"comment"`
		Code       string `json:"code"`
	} `json:"comments"`
}
