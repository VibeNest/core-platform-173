package entity

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrEmptyText    = errors.New("original text cannot be empty")
	ErrSameLanguage = errors.New("source and destination languages must be different")
)

// TranslationHistory — сущность, описывающая успешный перевод.
type TranslationHistory struct {
	ID          uuid.UUID `json:"id"`
	Source      string    `json:"source"`      // Код языка (en, ru)
	Destination string    `json:"destination"` // Код языка назначения
	Original    string    `json:"original"`    // Исходный текст
	Translation string    `json:"translation"` // Результат перевода
	CreatedAt   time.Time `json:"created_at"`
}

// Validate проверяет бизнес-правила сущности.
func (th *TranslationHistory) Validate() error {
	if strings.TrimSpace(th.Original) == "" {
		return ErrEmptyText
	}
	if th.Source == th.Destination {
		return ErrSameLanguage
	}
	return nil
}

// NewTranslationHistory — конструктор со встроенной нормализацией.
func NewTranslationHistory(src, dst, orig, trans string) *TranslationHistory {
	return &TranslationHistory{
		ID:          uuid.New(),
		Source:      strings.ToLower(strings.TrimSpace(src)),
		Destination: strings.ToLower(strings.TrimSpace(dst)),
		Original:    strings.TrimSpace(orig),
		Translation: strings.TrimSpace(trans),
		CreatedAt:   time.Now(),
	}
}
