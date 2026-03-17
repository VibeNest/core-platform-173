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

# Переменные для удобства (DRY)
DOCKER_COMPOSE = docker compose
# Указываем имя сервиса из docker-compose.yml
MIGRATE_SERVICE = migrate

.PHONY: migrate-create
migrate-create: ## Создать миграцию (использование: make migrate-create name=init_table)
	@if [ -z "$(name)" ]; then echo "Error: name is required. Use: make migrate-create name=my_migration"; exit 1; fi
	$(DOCKER_COMPOSE) run --rm --entrypoint migrate $(MIGRATE_SERVICE) create -ext sql -dir /migrations -seq $(name)

# В начале Makefile (где-то под импортом .env)
# Эти переменные подтянутся из твоего .env автоматически
PG_URL_DOCKER = postgres://user:password@postgres:5432/core_db?sslmode=disable

.PHONY: migrate-up
migrate-up: ## Применить все миграции (явная передача параметров)
	docker compose run --rm migrate -path=/migrations/ -database "$(PG_URL_DOCKER)" up

.PHONY: migrate-down
migrate-down: ## Откатить миграцию (явная передача параметров)
	docker compose run --rm migrate -path=/migrations/ -database "$(PG_URL_DOCKER)" down 1





