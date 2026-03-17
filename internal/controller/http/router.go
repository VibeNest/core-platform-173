package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	v1 "gitverse.ru/apavlov-systems/core-platform/internal/controller/http/v1"
	"gitverse.ru/apavlov-systems/core-platform/internal/usecase"
)

func NewRouter(handler *fiber.App, t usecase.Translation) {
	app := fiber.New()

	app.Use(logger.New())
	// Options
	handler.Use(logger.New())
	handler.Use(recover.New())

	// Health check
	handler.Get("/healthz", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

	// API v1
	v1.NewTranslationRoutes(handler.Group("/v1"), t)
}
