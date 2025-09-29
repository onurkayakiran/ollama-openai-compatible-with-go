package services

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"openai-compatible/config"
	"openai-compatible/models"
)

type OllamaService struct {
	config *config.Config
	client *http.Client
}

func NewOllamaService(cfg *config.Config) *OllamaService {
	return &OllamaService{
		config: cfg,
		client: &http.Client{
			Timeout: 300 * time.Second, // 5 minutes timeout for long responses
		},
	}
}

// Chat completion with Ollama
func (s *OllamaService) ChatCompletion(req *models.ChatCompletionRequest) (*models.ChatCompletionResponse, error) {
	// Use the model from request, fallback to config if empty
	modelName := req.Model
	if modelName == "" {
		modelName = s.config.OllamaModel
	}

	// Convert messages to ensure content is string for Ollama
	ollamaMessages := s.convertMessagesForOllama(req.Messages)

	ollamaReq := &models.OllamaChatRequest{
		Model:    modelName,
		Messages: ollamaMessages,
		Stream:   false,
		Options:  s.convertOptions(req),
	}

	jsonData, err := json.Marshal(ollamaReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := s.client.Post(s.config.OllamaURL+"/api/chat", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to make request to Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Ollama API error: %s", string(body))
	}

	var ollamaResp models.OllamaChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to decode Ollama response: %w", err)
	}

	// Convert to OpenAI format
	return &models.ChatCompletionResponse{
		ID:      generateID(),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   req.Model,
		Choices: []models.ChatCompletionChoice{
			{
				Index:        0,
				Message:      ollamaResp.Message,
				FinishReason: "stop",
			},
		},
		Usage: models.ChatCompletionUsage{
			PromptTokens:     estimateTokens(formatMessages(req.Messages)),
			CompletionTokens: estimateTokens(ollamaResp.Message.GetContentAsString()),
			TotalTokens:      estimateTokens(formatMessages(req.Messages)) + estimateTokens(ollamaResp.Message.GetContentAsString()),
		},
	}, nil
}

// Streaming chat completion
func (s *OllamaService) ChatCompletionStream(req *models.ChatCompletionRequest) (<-chan string, error) {
	// Use the model from request, fallback to config if empty
	modelName := req.Model
	if modelName == "" {
		modelName = s.config.OllamaModel
	}

	// Convert messages to ensure content is string for Ollama
	ollamaMessages := s.convertMessagesForOllama(req.Messages)

	ollamaReq := &models.OllamaChatRequest{
		Model:    modelName,
		Messages: ollamaMessages,
		Stream:   true,
		Options:  s.convertOptions(req),
	}

	jsonData, err := json.Marshal(ollamaReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := s.client.Post(s.config.OllamaURL+"/api/chat", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to make request to Ollama: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("Ollama API error: %s", string(body))
	}

	streamChan := make(chan string, 100)

	go func() {
		defer resp.Body.Close()
		defer close(streamChan)

		scanner := bufio.NewScanner(resp.Body)
		id := generateID()
		created := time.Now().Unix()

		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}

			var ollamaResp models.OllamaChatResponse
			if err := json.Unmarshal([]byte(line), &ollamaResp); err != nil {
				continue
			}

			// Convert to OpenAI streaming format
			streamResp := models.ChatCompletionStreamResponse{
				ID:      id,
				Object:  "chat.completion.chunk",
				Created: created,
				Model:   req.Model,
				Choices: []models.ChatCompletionStreamChoice{
					{
						Index: 0,
						Delta: models.ChatCompletionStreamDelta{
							Content: ollamaResp.Message.GetContentAsString(),
						},
					},
				},
			}

			if ollamaResp.Done {
				streamResp.Choices[0].FinishReason = stringPtr("stop")
			}

			jsonData, err := json.Marshal(streamResp)
			if err != nil {
				continue
			}

			streamChan <- "data: " + string(jsonData) + "\n\n"

			if ollamaResp.Done {
				streamChan <- "data: [DONE]\n\n"
				break
			}
		}
	}()

	return streamChan, nil
}

// Text completion
func (s *OllamaService) Completion(req *models.CompletionRequest) (*models.CompletionResponse, error) {
	prompt := ""
	switch p := req.Prompt.(type) {
	case string:
		prompt = p
	case []string:
		prompt = strings.Join(p, "\n")
	default:
		return nil, fmt.Errorf("unsupported prompt type")
	}

	ollamaReq := &models.OllamaGenerateRequest{
		Model:   s.config.OllamaModel,
		Prompt:  prompt,
		Stream:  false,
		Options: s.convertOptionsFromCompletion(req),
	}

	jsonData, err := json.Marshal(ollamaReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := s.client.Post(s.config.OllamaURL+"/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to make request to Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Ollama API error: %s", string(body))
	}

	var ollamaResp models.OllamaGenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to decode Ollama response: %w", err)
	}

	return &models.CompletionResponse{
		ID:      generateID(),
		Object:  "text_completion",
		Created: time.Now().Unix(),
		Model:   req.Model,
		Choices: []models.CompletionChoice{
			{
				Text:         ollamaResp.Response,
				Index:        0,
				FinishReason: "stop",
			},
		},
		Usage: models.CompletionUsage{
			PromptTokens:     estimateTokens(prompt),
			CompletionTokens: estimateTokens(ollamaResp.Response),
			TotalTokens:      estimateTokens(prompt) + estimateTokens(ollamaResp.Response),
		},
	}, nil
}

// Get available models
func (s *OllamaService) GetModels() (*models.ModelsResponse, error) {
	resp, err := s.client.Get(s.config.OllamaURL + "/api/tags")
	if err != nil {
		return nil, fmt.Errorf("failed to get models from Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Ollama API error: %s", string(body))
	}

	var ollamaResp models.OllamaModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to decode Ollama response: %w", err)
	}

	// Convert to OpenAI format
	openaiModels := make([]models.Model, 0, len(ollamaResp.Models))
	for _, model := range ollamaResp.Models {
		openaiModels = append(openaiModels, models.Model{
			ID:      model.Name,
			Object:  "model",
			Created: time.Now().Unix(),
			OwnedBy: "ollama",
		})
	}

	return &models.ModelsResponse{
		Object: "list",
		Data:   openaiModels,
	}, nil
}

// Helper functions
func (s *OllamaService) convertOptions(req *models.ChatCompletionRequest) *models.OllamaOptions {
	options := &models.OllamaOptions{}

	if req.Temperature != nil {
		options.Temperature = req.Temperature
	}
	if req.TopP != nil {
		options.TopP = req.TopP
	}
	if req.MaxTokens != nil {
		options.NumPredict = req.MaxTokens
	}
	if req.PresencePenalty != nil {
		options.PresencePenalty = req.PresencePenalty
	}
	if req.FrequencyPenalty != nil {
		options.FrequencyPenalty = req.FrequencyPenalty
	}

	// Handle stop sequences
	if req.Stop != nil {
		switch stop := req.Stop.(type) {
		case string:
			options.Stop = []string{stop}
		case []string:
			options.Stop = stop
		case []interface{}:
			stopStrings := make([]string, 0, len(stop))
			for _, s := range stop {
				if str, ok := s.(string); ok {
					stopStrings = append(stopStrings, str)
				}
			}
			options.Stop = stopStrings
		}
	}

	return options
}

func (s *OllamaService) convertOptionsFromCompletion(req *models.CompletionRequest) *models.OllamaOptions {
	options := &models.OllamaOptions{}

	if req.Temperature != nil {
		options.Temperature = req.Temperature
	}
	if req.TopP != nil {
		options.TopP = req.TopP
	}
	if req.MaxTokens != nil {
		options.NumPredict = req.MaxTokens
	}
	if req.PresencePenalty != nil {
		options.PresencePenalty = req.PresencePenalty
	}
	if req.FrequencyPenalty != nil {
		options.FrequencyPenalty = req.FrequencyPenalty
	}

	// Handle stop sequences
	if req.Stop != nil {
		switch stop := req.Stop.(type) {
		case string:
			options.Stop = []string{stop}
		case []string:
			options.Stop = stop
		case []interface{}:
			stopStrings := make([]string, 0, len(stop))
			for _, s := range stop {
				if str, ok := s.(string); ok {
					stopStrings = append(stopStrings, str)
				}
			}
			options.Stop = stopStrings
		}
	}

	return options
}

func generateID() string {
	return fmt.Sprintf("chatcmpl-%d", time.Now().UnixNano())
}

func estimateTokens(text string) int {
	// Simple token estimation (roughly 4 characters per token)
	return len(text) / 4
}

func formatMessages(messages []models.ChatMessage) string {
	var result strings.Builder
	for _, msg := range messages {
		result.WriteString(msg.Role + ": " + msg.GetContentAsString() + "\n")
	}
	return result.String()
}

// convertMessagesForOllama converts messages to ensure content is string for Ollama
func (s *OllamaService) convertMessagesForOllama(messages []models.ChatMessage) []models.ChatMessage {
	convertedMessages := make([]models.ChatMessage, len(messages))
	for i, msg := range messages {
		convertedMessages[i] = models.ChatMessage{
			Role:    msg.Role,
			Content: msg.GetContentAsString(), // Convert to string
			Name:    msg.Name,
		}
	}
	return convertedMessages
}

func stringPtr(s string) *string {
	return &s
}
