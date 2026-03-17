package v1

import (
	"github.com/gofiber/fiber/v2"
	"gitverse.ru/apavlov-systems/core-platform/internal/usecase"
)

type translationRoutes struct {
	t usecase.Translation
}

func NewTranslationRoutes(handler fiber.Router, t usecase.Translation) {
	r := &translationRoutes{t}

	h := handler.Group("/translation")
	{
		h.Post("/do", r.doTranslate)
		h.Get("/history", r.history)
	}
}

func (r *translationRoutes) doTranslate(c *fiber.Ctx) error {
	var request doTranslateRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	// Вызов UseCase
	res, err := r.t.Translate(c.Context(), request.Source, request.Destination, request.Text)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal error"})
	}

	return c.JSON(res)
}

func (r *translationRoutes) history(c *fiber.Ctx) error {
	// 1. Получаем сущности из UseCase
	items, err := r.t.History(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal error"})
	}

	// 2. Мапим []entity.TranslationHistory -> []historyResponse
	// Это важный архитектурный шаг: мы не отдаем ID или технические поля
	response := make([]historyResponse, len(items))
	for i, item := range items {
		response[i] = historyResponse{
			Source:      item.Source,
			Destination: item.Destination,
			Original:    item.Original,
			Translation: item.Translation,
		}
	}

	// 3. Отдаем уже подготовленный слайс DTO
	return c.JSON(fiber.Map{"history": response})
}
