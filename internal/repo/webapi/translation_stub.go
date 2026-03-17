package webapi

import (
	"context"
	"fmt"
	"strings"
)

// TranslationWebAPI — заглушка внешнего сервиса перевода.
// В будущем здесь может появиться http.Client для запросов к Google/DeepL.
type TranslationWebAPI struct {
	// Здесь можно добавить настройки, например, задержку имитации сети
}

// New — конструктор для создания инстанса «переводчика».
func New() *TranslationWebAPI {
	return &TranslationWebAPI{}
}

// Translate имитирует перевод текста. 
// Для наглядности он просто переводит текст в верхний регистр и добавляет метку языка.
func (w *TranslationWebAPI) Translate(ctx context.Context, source, destination, text string) (string, error) {
	// Имитируем работу: проверяем контекст (важно для Go-стиля)
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		// Логика «заглушки»
		translated := fmt.Sprintf("[%s->%s]: %s", 
			strings.ToUpper(source), 
			strings.ToUpper(destination), 
			strings.ToUpper(text),
		)
		return translated, nil
	}
}
