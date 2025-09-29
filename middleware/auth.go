package middleware

import (
	"strings"

	"openai-compatible/config"
	"openai-compatible/models"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(models.ErrorResponse{
				Error: models.ErrorDetail{
					Message: "Authorization header is required",
					Type:    "invalid_request_error",
					Code:    "missing_authorization",
				},
			})
		}

		// Check if it starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(401).JSON(models.ErrorResponse{
				Error: models.ErrorDetail{
					Message: "Authorization header must start with 'Bearer '",
					Type:    "invalid_request_error",
					Code:    "invalid_authorization_format",
				},
			})
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token != cfg.APIKey {
			return c.Status(401).JSON(models.ErrorResponse{
				Error: models.ErrorDetail{
					Message: "Invalid API key",
					Type:    "invalid_request_error",
					Code:    "invalid_api_key",
				},
			})
		}

		return c.Next()
	}
}
