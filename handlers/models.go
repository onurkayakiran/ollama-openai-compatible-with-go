package handlers

import (
	"log"

	"openai-compatible/models"
	"openai-compatible/services"

	"github.com/gofiber/fiber/v2"
)

type ModelsHandler struct {
	ollamaService *services.OllamaService
}

func NewModelsHandler(ollamaService *services.OllamaService) *ModelsHandler {
	return &ModelsHandler{
		ollamaService: ollamaService,
	}
}

func (h *ModelsHandler) GetModels(c *fiber.Ctx) error {
	resp, err := h.ollamaService.GetModels()
	if err != nil {
		log.Printf("Error getting models: %v", err)
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
