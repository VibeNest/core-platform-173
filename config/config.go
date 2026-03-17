package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type (
	// Config — корневая структура конфигурации всего приложения.
	// Использует вложенные структуры для логического разделения настроек.
	Config struct {
		App     App
		HTTP    HTTP
		Log     Log
		PG      PG
		GRPC    GRPC
		RMQ     RMQ
		NATS    NATS
		Swagger SWAGGER
		Metrics METRICS
	}

	// App содержит метаданные приложения, необходимые для логов и мониторинга.
	App struct {
		Name    string `env:"APP_NAME,required"`
		Version string `env:"APP_VERSION,required"`
		Env     string `env:"APP_ENV" envDefault:"dev"`
	}

	// HTTP описывает настройки веб-сервера.
	HTTP struct {
		Port string `env:"HTTP_PORT,required"`
		// UsePreforkMode позволяет включить режим Fiber Prefork (актуально для высоконагруженных систем).
		UsePreforkMode bool `env:"HTTP_USE_PREFORK_MODE" envDefault:"false"`
	}

	// Log задает уровень детализации логов (debug, info, warn, error).
	Log struct {
		Level string `env:"LOG_LEVEL,required"`
	}

	// PG содержит параметры подключения к PostgreSQL.
	PG struct {
		PoolMax int    `env:"PG_POOL_MAX,required"` // Максимальное количество соединений в пуле.
		URL     string `env:"PG_URL,required"`      // DSN строка подключения.
	}

	// GRPC описывает порт для межсервисного взаимодействия через gRPC.
	GRPC struct {
		Port string `env:"GRPC_PORT,required"`
	}

	// RMQ содержит настройки для RabbitMQ (AMQP).
	RMQ struct {
		ServerExchange string `env:"RMQ_RPC_SERVER,required"`
		ClientExchange string `env:"RMQ_RPC_CLIENT,required"`
		URL            string `env:"RMQ_URL,required"`
	}

	// NATS содержит настройки для работы с NATS JetStream или Core NATS.
	NATS struct {
		ServerExchange string `env:"NATS_RPC_SERVER,required"`
		URL            string `env:"NATS_URL,required"`
	}

	METRICS struct {
		Enabled bool `env:"METRICS_ENABLED" envDefault:"true"`
	}

	SWAGGER struct {
		Enabled bool `env:"SWAGGER_ENABLED" envDefault:"false"`
	}
)

// NewConfig инициализирует структуру Config, считывая данные из переменных окружения.
// Если обязательная переменная (tag: required) отсутствует, вернет ошибку.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	// Парсинг переменных окружения напрямую в структуру.
	// env.Parse поддерживает вложенные структуры автоматически.
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config parse error: %w", err)
	}

	return cfg, nil
}
