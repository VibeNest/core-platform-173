# Подгружаем переменные из .env, если он существует.
# Если нет — берем из .env.example, чтобы проект не падал при первом запуске.
ifneq (,$(wildcard ./.env))
    include .env
    export
else ifneq (,$(wildcard ./.env.example))
    include .env.example
    export
endif

# Указываем цель по умолчанию (если просто написать make в консоли)
.DEFAULT_GOAL := help

.PHONY: deps
deps: ## Загрузка и проверка зависимостей Go
	go mod tidy
	go mod verify

.PHONY: run
run: ## Локальный запуск приложения (с подгрузкой ENV)
	go run ./cmd/app/main.go

.PHONY: help
help: ## Список доступных команд с описанием
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: compose-up
compose-up: ## Запустить инфраструктуру (Postgres) в Docker
	docker compose up -d

.PHONY: compose-down
compose-down: ## Остановить все контейнеры
	docker compose down
