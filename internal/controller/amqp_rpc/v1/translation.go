package v1

import (
	"context"
	"encoding/json" // Добавь импорт для работы с JSON

	"gitverse.ru/apavlov-systems/core-platform/internal/usecase"
)

type translationRoutes struct {
	t usecase.Translation
}

// Вспомогательная структура для парсинга запроса из RabbitMQ
type translateRequest struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Original    string `json:"original"`
}

func (r *translationRoutes) translate(ctx context.Context, data []byte) (interface{}, error) {
	// 1. Распаковываем входящие байты в структуру
	var req translateRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}

	// 2. Вызываем бизнес-логику (UseCase)
	// Она переведет текст и сохранит его в базу
	res, err := r.t.Translate(ctx, req.Source, req.Destination, req.Original)
	if err != nil {
		return nil, err
	}

	// 3. Возвращаем результат (теперь это не nil, а объект из базы)
	return res, nil
}

func (r *translationRoutes) history(ctx context.Context, data []byte) (interface{}, error) {
	history, err := r.t.History(ctx)
	if err != nil {
		return nil, err
	}
	return history, nil
}
