package usecase

import (
	"context"
	"strings"

	"gitverse.ru/apavlov-systems/core-platform/internal/entity"
)

// MockRepo — временное хранилище в оперативной памяти
type MockRepo struct {
	data []entity.TranslationHistory
}

func (m *MockRepo) StoreHistory(ctx context.Context, t entity.TranslationHistory) error {
	m.data = append(m.data, t)
	return nil
}

func (m *MockRepo) GetHistory(ctx context.Context) ([]entity.TranslationHistory, error) {
	return m.data, nil
}

// MockTranslator — имитация внешнего сервиса
type MockTranslator struct{}

func (m *MockTranslator) Translate(ctx context.Context, src, dst, text string) (string, error) {
	// Имитируем «перевод», просто переводя текст в верхний регистр
	return strings.ToUpper(text), nil
}
