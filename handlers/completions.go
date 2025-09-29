package handlers

import (
	"log"

	"openai-compatible/models"
	"openai-compatible/services"

	"github.com/gofiber/fiber/v2"
)

type CompletionsHandler struct {
	ollamaService *services.OllamaService
}

func NewCompletionsHandler(ollamaService *services.OllamaService) *CompletionsHandler {
	return &CompletionsHandler{
		ollamaService: ollamaService,
	}
}

func (h *CompletionsHandler) Completions(c *fiber.Ctx) error {
	var req models.CompletionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{
			Error: models.ErrorDetail{
				Message: "Invalid request body",
				Type:    "invalid_request_error",
				Code:    "invalid_json",
			},
		})
	}

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

	if req.Prompt == nil {
		return c.Status(400).JSON(models.ErrorResponse{
			Error: models.ErrorDetail{
				Message: "Prompt is required",
				Type:    "invalid_request_error",
				Code:    "missing_prompt",
			},
		})
	}

	// Handle request
	resp, err := h.ollamaService.Completion(&req)
	if err != nil {
		log.Printf("Error in completion: %v", err)
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
