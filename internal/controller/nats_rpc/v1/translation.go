package v1

import (
	"context"
	"encoding/json"

	"gitverse.ru/apavlov-systems/core-platform/internal/usecase"
)

type translationRoutes struct {
	t usecase.Translation
}

func newTranslationRoutes(t usecase.Translation) *translationRoutes {
	return &translationRoutes{t}
}

// Request/Response структуры для JSON
type translateRequest struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Original    string `json:"original"`
}

func (r *translationRoutes) translate(ctx context.Context, data []byte) (interface{}, error) {
	var req translateRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err // В будущем тут будет кастомная ошибка
	}

	// Вызов бизнес-логики
	res, err := r.t.Translate(ctx, req.Source, req.Destination, req.Original)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *translationRoutes) history(ctx context.Context, data []byte) (interface{}, error) {
	// Вызываем usecase для получения истории
	history, err := r.t.History(ctx)
	if err != nil {
		return nil, err
	}

	return history, nil
}
