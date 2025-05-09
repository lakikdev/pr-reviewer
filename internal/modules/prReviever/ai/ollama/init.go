package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pr-reviewer/internal/config"
	"pr-reviewer/internal/modules/prReviever/ai/model"
	"strings"
)

const SystemContent = `
You are a senior software engineer performing a focused pull request review.

Only analyze the lines **added** in the diff. Ignore all removed lines and any unchanged lines.
Do not comment on deleted lines or code that no longer exists.
Ignore non-functional strings like documentation, comments, or embedded prompts unless they affect runtime behavior.
You will be provided with Diff section with all the changes made in the pull request.
symbol '+' indicates the lines added in the pull request.
symbol '-' indicates the lines removed in the pull request.

Use Diff to identify the changes made in the pull request and than use full code section to understand the context of the changes.
Your task is to provide **technical code review comments** for newly added lines only. Focus on logic, bugs, performance, or maintainability.
Avoid comments about personal style, documentation format, or text prompt content.

Do not include any comments if there are no meaningful issues. Do not review removed lines or string literals unless they affect execution.


Your goals:
- Identify any issues introduced in the modified lines.
- Suggest improvements related to clarity, readability, or code quality.
- Ignore unrelated parts of the file or project.

Respond with comments using JSON array format:

[
	{ 
		"line_number": <line number of the code with issue>, 
		"code": "<the code line with issue>",
 		"comment": "<your very detailed comment, with explaining the issue and suggestion of how to fix it>" 
	}
]

Example:
[
  {
    "line_number": 42,
	"code": "<the code line with issue>",
    "comment": "You should log the error before returning to aid in debugging."
  }
]


Here is the diff:
`

type OllamaClient struct {
	BaseURL string
	Model   string
}

func NewClient() *OllamaClient {
	return &OllamaClient{
		BaseURL: *config.OllamaBaseURL,
		Model:   *config.OllamaModel,
	}
}

type OllamaChatRequest struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
	Stream bool              `json:"stream"`
	Format *StructuredFormat `json:"format,omitempty"`
}

type GenerateResponse struct {
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
}

type ItemProperty struct {
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
}

type FormatProperty struct {
	Type        string       `json:"type"`
	Description string       `json:"description,omitempty"`
	Enum        []string     `json:"enum,omitempty"`
	Items       ItemProperty `json:"items,omitempty"`
}

type StructuredFormat struct {
	Type        string                      `json:"type"`
	Description string                      `json:"description,omitempty"`
	Properties  map[string]StructuredFormat `json:"properties"`
	Required    []string                    `json:"required,omitempty"`
	Items       map[string]ItemProperty     `json:"items,omitempty"`
}

func (o *OllamaClient) Send(prompt string) (string, error) {

	reqBody := OllamaChatRequest{
		Model: o.Model,
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{
				Role:    "system",
				Content: SystemContent,
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	//fmt.Printf("Request JSON: %s\n", string(jsonData)) // Debugging line
	// Send the HTTP POST request
	resp, err := http.Post(fmt.Sprint(o.BaseURL, "/api/chat"), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	//fmt.Printf("Response JSON: %s\n", string(body)) // Debugging line

	// Unmarshal JSON
	var genResp GenerateResponse
	err = json.Unmarshal(body, &genResp)
	if err != nil {
		fmt.Println("Raw Response:", string(body)) // fallback
		return "", err
	}

	return genResp.Message.Content, nil
}

func (o *OllamaClient) ParseResponse(resp string) model.StructResponse {
	trimmedResp := resp[strings.Index(resp, "[") : strings.Index(resp, "]")+1]
	// fmt.Printf("Trimmed Response: %s\n", trimmedResp) // Debugging line

	parsedResponse := model.StructResponse{}
	err := json.Unmarshal([]byte(trimmedResp), &parsedResponse.Comments)
	if err != nil {
		fmt.Println("Error parsing response:", err)
		fmt.Println("Raw Response:", resp) // fallback
		return model.StructResponse{}
	}

	return parsedResponse

}
