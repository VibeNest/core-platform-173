package v1

import (
	"gitverse.ru/apavlov-systems/core-platform/internal/usecase"
	"gitverse.ru/apavlov-systems/core-platform/pkg/amqprpc"
	"gitverse.ru/apavlov-systems/core-platform/pkg/logger"
)

func RegisterRoutes(server *amqprpc.Server, t usecase.Translation, l *logger.Logger) {
	routes := &translationRoutes{t, l}

	// Регистрируем те же методы, что и в NATS
	server.Register("v1.translation.translate", routes.translate)
	server.Register("v1.translation.history", routes.history)
}
