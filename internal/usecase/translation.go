package usecase

import (
	"context"
	"errors"
	"fmt"

	"gitverse.ru/apavlov-systems/core-platform/internal/entity"
	"gitverse.ru/apavlov-systems/core-platform/pkg/logger"
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
	webAPI TranslationWebAPI
	l      *logger.Logger
}

// New — конструктор с Dependency Injection (DI).
func New(r TranslationRepo, w TranslationWebAPI, l *logger.Logger) *TranslationUseCase {
	return &TranslationUseCase{
		repo:   r,
		webAPI: w,
		l: l,
	}
}

func (uc *TranslationUseCase) Translate(ctx context.Context, src, dst, text string) (entity.TranslationHistory, error) {
	uc.l.Debug("UseCase: translating %s to %s", src, dst)
	// 1. Получаем СТРОКУ от внешнего API
	translatedText, err := uc.webAPI.Translate(ctx, src, dst, text)
	if err != nil {
		return entity.TranslationHistory{}, err
	}

	// 2. САМИ создаем структуру истории
	history := entity.TranslationHistory{
		Source:      src,
		Destination: dst,
		Original:    text,
		Translation: translatedText, // кладем строку сюда
	}

	// 3. Сохраняем в базу через репозиторий
	_ = uc.repo.StoreHistory(ctx, history)

	return history, nil
}

// History — получение списка последних переводов.
func (uc *TranslationUseCase) History(ctx context.Context) ([]entity.TranslationHistory, error) {
	histories, err := uc.repo.GetHistory(ctx)
	if err != nil {
		return nil, fmt.Errorf("usecase - History - repo.GetHistory: %w", ErrInternal)
	}

	return histories, nil
}
