package v1

import (
	"context"

	"gitverse.ru/apavlov-systems/core-platform/internal/usecase"
)

type translationRoutes struct {
	t usecase.Translation
}

func (r *translationRoutes) translate(ctx context.Context, data []byte) (interface{}, error) {
	// 1. Unmarshal JSON -> req
	// 2. res, err := r.t.Translate(...)
	return nil, nil
}

func (r *translationRoutes) history(ctx context.Context, data []byte) (interface{}, error) {
	// Вызываем бизнес-логику получения истории
	history, err := r.t.History(ctx)
	if err != nil {
		return nil, err
	}

	// Возвращаем результат (движок сам превратит его в JSON)
	return history, nil
}
