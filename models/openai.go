package models

import (
	"fmt"
	"strings"
)

// Chat Completions Request
type ChatCompletionRequest struct {
	Model            string                 `json:"model"`
	Messages         []ChatMessage          `json:"messages"`
	MaxTokens        *int                   `json:"max_tokens,omitempty"`
	Temperature      *float64               `json:"temperature,omitempty"`
	TopP             *float64               `json:"top_p,omitempty"`
	N                *int                   `json:"n,omitempty"`
	Stream           *bool                  `json:"stream,omitempty"`
	Stop             interface{}            `json:"stop,omitempty"`
	PresencePenalty  *float64               `json:"presence_penalty,omitempty"`
	FrequencyPenalty *float64               `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]interface{} `json:"logit_bias,omitempty"`
	User             string                 `json:"user,omitempty"`
}

type ChatMessage struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
	Name    string      `json:"name,omitempty"`
}

// GetContentAsString returns content as string, handling both string and array formats
func (m *ChatMessage) GetContentAsString() string {
	switch content := m.Content.(type) {
	case string:
		return content
	case []interface{}:
		var parts []string
		for _, part := range content {
			if str, ok := part.(string); ok {
				parts = append(parts, str)
			} else if partMap, ok := part.(map[string]interface{}); ok {
				// Handle structured content like {"type": "text", "text": "content"}
				if text, exists := partMap["text"]; exists {
					if textStr, ok := text.(string); ok {
						parts = append(parts, textStr)
					}
				}
			}
		}
		return strings.Join(parts, " ")
	case []string:
		return strings.Join(content, " ")
	default:
		return fmt.Sprintf("%v", content)
	}
}

// Chat Completions Response
type ChatCompletionResponse struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int64                  `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChatCompletionChoice `json:"choices"`
	Usage   ChatCompletionUsage    `json:"usage"`
}

type ChatCompletionChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

type ChatCompletionUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Streaming Response
type ChatCompletionStreamResponse struct {
	ID      string                       `json:"id"`
	Object  string                       `json:"object"`
	Created int64                        `json:"created"`
	Model   string                       `json:"model"`
	Choices []ChatCompletionStreamChoice `json:"choices"`
}

type ChatCompletionStreamChoice struct {
	Index        int                       `json:"index"`
	Delta        ChatCompletionStreamDelta `json:"delta"`
	FinishReason *string                   `json:"finish_reason"`
}

type ChatCompletionStreamDelta struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

// Text Completions Request
type CompletionRequest struct {
	Model            string      `json:"model"`
	Prompt           interface{} `json:"prompt"`
	MaxTokens        *int        `json:"max_tokens,omitempty"`
	Temperature      *float64    `json:"temperature,omitempty"`
	TopP             *float64    `json:"top_p,omitempty"`
	N                *int        `json:"n,omitempty"`
	Stream           *bool       `json:"stream,omitempty"`
	Logprobs         *int        `json:"logprobs,omitempty"`
	Echo             *bool       `json:"echo,omitempty"`
	Stop             interface{} `json:"stop,omitempty"`
	PresencePenalty  *float64    `json:"presence_penalty,omitempty"`
	FrequencyPenalty *float64    `json:"frequency_penalty,omitempty"`
	BestOf           *int        `json:"best_of,omitempty"`
	User             string      `json:"user,omitempty"`
}

// Text Completions Response
type CompletionResponse struct {
	ID      string             `json:"id"`
	Object  string             `json:"object"`
	Created int64              `json:"created"`
	Model   string             `json:"model"`
	Choices []CompletionChoice `json:"choices"`
	Usage   CompletionUsage    `json:"usage"`
}

type CompletionChoice struct {
	Text         string `json:"text"`
	Index        int    `json:"index"`
	Logprobs     *int   `json:"logprobs"`
	FinishReason string `json:"finish_reason"`
}

type CompletionUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Models Response
type ModelsResponse struct {
	Object string  `json:"object"`
	Data   []Model `json:"data"`
}

type Model struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

// Error Response
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code,omitempty"`
}
