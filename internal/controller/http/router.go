package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"gitverse.ru/apavlov-systems/core-platform/config"
	v1 "gitverse.ru/apavlov-systems/core-platform/internal/controller/http/v1"
	"gitverse.ru/apavlov-systems/core-platform/internal/usecase"
)

func NewRouter(handler *fiber.App, t usecase.Translation, cfg *config.Config) {
	app := fiber.New()

	app.Use(logger.New())
	// Options
	handler.Use(logger.New())
	handler.Use(recover.New())

	// Health check
	handler.Get("/healthz", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

	// Metrics (Prometheus) - Добавим чуть позже через библиотеку
	if cfg.Metrics.Enabled {
		// handler.Get("/metrics", monitor.New())
	}

	if cfg.Swagger.Enabled {
		// Здесь позже добавим подключение самого Swagger-обработчика
		handler.Get("/swagger/*", func(c *fiber.Ctx) error { return c.SendString("Swagger is coming soon...") })
	}

	// API v1
	v1.NewTranslationRoutes(handler.Group("/v1"), t)
}
