package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/nats-io/nats.go"
	amqp "github.com/rabbitmq/amqp091-go"
	"gitverse.ru/apavlov-systems/core-platform/config"
	amqp_ctrl "gitverse.ru/apavlov-systems/core-platform/internal/controller/amqp_rpc/v1"
	grpc_ctrl "gitverse.ru/apavlov-systems/core-platform/internal/controller/grpc/v1"
	"gitverse.ru/apavlov-systems/core-platform/internal/controller/http"
	nats_ctrl "gitverse.ru/apavlov-systems/core-platform/internal/controller/nats_rpc/v1"
	"gitverse.ru/apavlov-systems/core-platform/internal/repo/persistent"
	"gitverse.ru/apavlov-systems/core-platform/internal/repo/postgres"
	"gitverse.ru/apavlov-systems/core-platform/internal/repo/webapi"
	"gitverse.ru/apavlov-systems/core-platform/internal/usecase"
	"gitverse.ru/apavlov-systems/core-platform/pkg/amqprpc"
	"gitverse.ru/apavlov-systems/core-platform/pkg/httpserver"
	"gitverse.ru/apavlov-systems/core-platform/pkg/natsrpc"
	"google.golang.org/grpc"
)

// Run — “продолжение main”: тут будет DI и запуск серверов.
func Run(cfg *config.Config) {
	// 1. Инициализация Postgres (Infrastructure)
	// Мы используем нашу обертку из pkg/postgres, которая умеет ждать базу
	pg, err := postgres.New(cfg.PG.URL) // Добавь настройки PoolMax если реализовал в pkg
	if err != nil {
		log.Fatalf("app - Run - postgres.New: %v", err)
	}
	defer pg.Close()

	// Инициализация Репозиториев (Adapters)
	// Передаем подключение к базе в репозиторий истории
	historyRepo := persistent.New(pg)

	// Инициализация Внешних API (Adapters)
	// Наша заглушка-переводчик
	translator := webapi.New()

	// Инициализация UseCase (Бизнес-логика - ОДНА для всех)
	translationUseCase := usecase.New(historyRepo, translator)

	// --- Инициализация Транспортов ---

	// HTTP Server (REST)
	// 1. Создаем объект приложения (Fiber), который просит NewRouter
	app := fiber.New()
	http.NewRouter(app, translationUseCase, cfg)

	// Используем adaptor.FiberApp, чтобы превратить app в http.Handler
	// И используем "=", так как httpServer уже мог быть объявлен выше (ошибка NoNewVar)
	httpServer := httpserver.New(adaptor.FiberApp(app), httpserver.Port(cfg.HTTP.Port))

	// 4. gRPC Server
	gRPCServer := grpc.NewServer()

	grpc_ctrl.RegisterRoutes(gRPCServer, translationUseCase)
	// 5. NATS RPC
	nc, err := nats.Connect(cfg.NATS.URL)
	if err != nil {
		log.Fatalf("app - Run - nats.Connect: %v", err)
	}
	defer nc.Close()
	natsServer := natsrpc.NewServer(nc)
	nats_ctrl.RegisterRoutes(natsServer, translationUseCase)

	// 6. AMQP (RabbitMQ) RPC
	rmqConn, err := amqp.Dial(cfg.RMQ.URL)
	if err != nil {
		log.Fatalf("app - Run - amqp.Dial: %v", err)
	}
	defer rmqConn.Close()
	rmqChan, _ := rmqConn.Channel()
	rmqServer := amqprpc.NewServer(rmqConn, rmqChan)
	amqp_ctrl.RegisterRoutes(rmqServer, translationUseCase)

	// --- Запуск серверов ---

	notify := make(chan error, 1)

	// gRPC
	go func() {
		listener, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
		if err != nil {
			notify <- fmt.Errorf("grpc listen: %w", err)
			return
		}
		log.Printf("app - Run - gRPC server listing on %s", cfg.GRPC.Port)
		notify <- gRPCServer.Serve(listener)
	}()

	// --- Ожидание завершения (Graceful Shutdown) ---

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Printf("app - Run - signal: " + s.String())
	case err := <-notify:
		log.Printf("app - Run - server error: %v", err)
	}

	// Порядок остановки
	httpServer.Shutdown()
	gRPCServer.GracefulStop()
	log.Printf("app - Run - stopped")
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
