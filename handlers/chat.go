package handlers

import (
	"bufio"
	"fmt"
	"log"

	"openai-compatible/models"
	"openai-compatible/services"

	"github.com/gofiber/fiber/v2"
)

type ChatHandler struct {
	ollamaService *services.OllamaService
}

func NewChatHandler(ollamaService *services.OllamaService) *ChatHandler {
	return &ChatHandler{
		ollamaService: ollamaService,
	}
}

func (h *ChatHandler) ChatCompletions(c *fiber.Ctx) error {
	// Log the raw request body for debugging
	body := c.Body()
	log.Printf("Received request body: %s", string(body))
	log.Printf("Content-Type: %s", c.Get("Content-Type"))
	log.Printf("Authorization: %s", c.Get("Authorization"))

	var req models.ChatCompletionRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("JSON parsing error: %v", err)
		return c.Status(400).JSON(models.ErrorResponse{
			Error: models.ErrorDetail{
				Message: fmt.Sprintf("Invalid request body: %v", err),
				Type:    "invalid_request_error",
				Code:    "invalid_json",
			},
		})
	}

	log.Printf("Parsed request: Model=%s, Messages=%d", req.Model, len(req.Messages))

	// Validate required fields
	if req.Model == "" {
		return c.Status(400).JSON(models.ErrorResponse{
			Error: models.ErrorDetail{
				Message: "Model is required",
				Type:    "invalid_request_error",
				Code:    "missing_model",
			},
		})
	}

	if len(req.Messages) == 0 {
		return c.Status(400).JSON(models.ErrorResponse{
			Error: models.ErrorDetail{
				Message: "Messages are required",
				Type:    "invalid_request_error",
				Code:    "missing_messages",
			},
		})
	}

	// Check if streaming is requested
	if req.Stream != nil && *req.Stream {
		return h.handleStreamingChat(c, &req)
	}

	// Handle non-streaming request
	resp, err := h.ollamaService.ChatCompletion(&req)
	if err != nil {
		log.Printf("Error in chat completion: %v", err)
		return c.Status(500).JSON(models.ErrorResponse{
			Error: models.ErrorDetail{
				Message: "Internal server error",
				Type:    "internal_error",
				Code:    "ollama_error",
			},
		})
	}

	return c.JSON(resp)
}

func (h *ChatHandler) handleStreamingChat(c *fiber.Ctx, req *models.ChatCompletionRequest) error {
	// Set headers for Server-Sent Events
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "Cache-Control")

	streamChan, err := h.ollamaService.ChatCompletionStream(req)
	if err != nil {
		log.Printf("Error in streaming chat completion: %v", err)
		return c.Status(500).JSON(models.ErrorResponse{
			Error: models.ErrorDetail{
				Message: "Internal server error",
				Type:    "internal_error",
				Code:    "ollama_stream_error",
			},
		})
	}

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic in stream writer: %v", r)
			}
		}()

		for data := range streamChan {
			if _, err := w.WriteString(data); err != nil {
				log.Printf("Error writing stream data: %v", err)
				break
			}
			if err := w.Flush(); err != nil {
				log.Printf("Error flushing stream: %v", err)
				break
			}
		}
	})

	return nil
}
