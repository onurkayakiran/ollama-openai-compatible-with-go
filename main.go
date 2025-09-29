package main

import (
	"log"

	"openai-compatible/config"
	"openai-compatible/handlers"
	"openai-compatible/middleware"
	"openai-compatible/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ServerHeader: "OpenAI-Compatible-API",
		AppName:      "OpenAI Compatible API v1.0.0",
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Initialize services
	ollamaService := services.NewOllamaService(cfg)

	// Initialize handlers
	chatHandler := handlers.NewChatHandler(ollamaService)
	completionsHandler := handlers.NewCompletionsHandler(ollamaService)
	modelsHandler := handlers.NewModelsHandler(ollamaService)

	// Health check endpoint (no auth required)
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "OpenAI Compatible API is running",
		})
	})

	// API routes with authentication
	api := app.Group("/v1", middleware.AuthMiddleware(cfg))

	// OpenAI compatible endpoints
	api.Post("/chat/completions", chatHandler.ChatCompletions)
	api.Post("/completions", completionsHandler.Completions)
	api.Get("/models", modelsHandler.GetModels)

	// Start server
	log.Printf("Starting server on port %s", cfg.Port)
	log.Printf("API Key: %s", cfg.APIKey)
	log.Printf("Ollama URL: %s", cfg.OllamaURL)
	log.Printf("Ollama Model: %s", cfg.OllamaModel)

	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
