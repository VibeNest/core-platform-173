package v1

import (
	"gitverse.ru/apavlov-systems/core-platform/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func RegisterRoutes(server *grpc.Server, t usecase.Translation) {
	// Включаем "просмотр" методов для grpcurl
	reflection.Register(server)

	// Создаем наш хендлер и регистрируем его в gRPC сервере
	handler := NewTranslationRoutes(t)
	RegisterTranslationServer(server, handler)
}
