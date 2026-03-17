package usecase

import (
	"context"
	"errors"
	"fmt"

	"gitverse.ru/apavlov-systems/core-platform/internal/entity"
)

var (
	// ErrInternal — общая ошибка сервера, когда что-то пошло не так во внешних системах.
	ErrInternal = errors.New("internal error")
	// ErrInvalidInput — ошибка валидации входных данных.
	ErrInvalidInput = errors.New("invalid input")
)

// TranslationUseCase — основной узел бизнес-логики.
type TranslationUseCase struct {
	repo   TranslationRepo
	webAPI Translation
}

// New — конструктор с Dependency Injection (DI).
func New(r TranslationRepo, w Translation) *TranslationUseCase {
	return &TranslationUseCase{
		repo:   r,
		webAPI: w,
	}
}

// Translate — главный бизнес-сценарий: перевести текст и сохранить историю.
func (uc *TranslationUseCase) Translate(ctx context.Context, src, dst, text string) (entity.TranslationHistory, error) {
	// 1. Создаем черновик сущности (здесь сработает нормализация строк)
	history := entity.NewTranslationHistory(src, dst, text, "")

	// 2. Бизнес-валидация сущности
	if err := history.Validate(); err != nil {
		return entity.TranslationHistory{}, fmt.Errorf("usecase - Translate - validate: %w", ErrInvalidInput)
	}

	// 3. Вызов внешнего переводчика
	translated, err := uc.webAPI.Translate(ctx, history.Source, history.Destination, history.Original)
	if err != nil {
		return entity.TranslationHistory{}, fmt.Errorf("usecase - Translate - webAPI.Translate: %w", ErrInternal)
	}

	// 4. Обновляем сущность результатом перевода
	history.Translation = translated

	// 5. Сохранение в репозиторий (базу данных)
	err = uc.repo.StoreHistory(ctx, *history)
	if err != nil {
		return entity.TranslationHistory{}, fmt.Errorf("usecase - Translate - repo.StoreHistory: %w", ErrInternal)
	}

	return *history, nil
}

// History — получение списка последних переводов.
func (uc *TranslationUseCase) History(ctx context.Context) ([]entity.TranslationHistory, error) {
	histories, err := uc.repo.GetHistory(ctx)
	if err != nil {
		return nil, fmt.Errorf("usecase - History - repo.GetHistory: %w", ErrInternal)
	}

	return histories, nil
}
