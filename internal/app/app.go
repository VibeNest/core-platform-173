package app

import (
	"context"
	"log"
	"time"

	"gitverse.ru/apavlov-systems/core-platform/config"
	"gitverse.ru/apavlov-systems/core-platform/internal/repo/persistent"
	"gitverse.ru/apavlov-systems/core-platform/internal/repo/postgres"
	"gitverse.ru/apavlov-systems/core-platform/internal/repo/webapi"
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

	// 1. Инициализация Postgres (Infrastructure)
	// Мы используем нашу обертку из pkg/postgres, которая умеет ждать базу
	pg, err := postgres.New(cfg.PG.URL) // Добавь настройки PoolMax если реализовал в pkg
	if err != nil {
		log.Fatalf("app - Run - postgres.New: %v", err)
	}
	defer pg.Close()

	// 2. Инициализация Репозиториев (Adapters)
	// Передаем подключение к базе в репозиторий истории
	historyRepo := persistent.New(pg)

	// 3. Инициализация Внешних API (Adapters)
	// Наша заглушка-переводчик
	translator := webapi.New()

	// 4. Инициализация UseCase (Business Logic / Core)
	// Соединяем "руки" (repo/webapi) с "мозгом" (usecase)
	translationUseCase := usecase.New(historyRepo, translator)

	// 5. Проверка работоспособности (Health Check)
	// В режиме dev сделаем тестовый вызов, чтобы убедиться, что всё связано верно
	if cfg.App.Env == "dev" {
		testUseCase(translationUseCase)
	}

	log.Printf("Приложение %s готово и ожидает транспортный уровень...", cfg.App.Name)

}

// Вспомогательная функция для "прогрева" системы
func testUseCase(uc *usecase.TranslationUseCase) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := uc.Translate(ctx, "en", "ru", "Глубокое погружение в чистую архитектуру")
	if err != nil {
		log.Printf("[TEST] Сбой выполнения сценария UseCase: %v", err)
		return
	}

	log.Printf("[TEST] Успешно! Оригинал: %s | Перевод: %s", res.Original, res.Translation)
}
