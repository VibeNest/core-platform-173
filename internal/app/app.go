package app

import (
	"context"
	"fmt"
	"log"

	"gitverse.ru/apavlov-systems/core-platform/config"
	"gitverse.ru/apavlov-systems/core-platform/internal/usecase"
)

// Run — “продолжение main”: тут будет DI и запуск серверов.
func Run(cfg *config.Config) {
	_ = cfg
	if cfg.App.Env == "dev" {
		log.Println("🛠 Запущено в режиме РАЗРАБОТКИ.")
	} else {
		log.Println("🚀 Запущено в режиме ПРОДАКШЕН. Максимальная производительность.")
	}

	// 1. Инициализируем заглушки (вместо реальных БД и API)
	repo := &usecase.MockRepo{}
	translator := &usecase.MockTranslator{}

	// 2. Создаем UseCase (Dependency Injection)
	uc := usecase.New(repo, translator)

	// 3. Пробуем выполнить бизнес-сценарий
	ctx := context.Background()
	res, err := uc.Translate(ctx, "en", "ru", "hello world")
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	fmt.Printf("Результат: %s -> %s\n", res.Original, res.Translation)

	// 4. Проверяем историю
	history, _ := uc.History(ctx)
	fmt.Printf("Записей в истории: %d\n", len(history))
}
