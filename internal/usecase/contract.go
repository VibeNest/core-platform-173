package usecase

import (
	"context"

	"gitverse.ru/apavlov-systems/core-platform/internal/entity"
)

type (
	// Translation — интерфейс внешнего сервиса перевода (Web API, Mock или Stub).
	Translation interface {
		Translate(ctx context.Context, source, destination, text string) (string, error)
		History(ctx context.Context) ([]entity.TranslationHistory, error)
	}

	// TranslationRepo — интерфейс для работы с базой данных (PostgreSQL, MongoDB).
	TranslationRepo interface {
		StoreHistory(ctx context.Context, t entity.TranslationHistory) error
		GetHistory(ctx context.Context) ([]entity.TranslationHistory, error)
	}
)
