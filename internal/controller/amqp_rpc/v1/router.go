package v1

import (
	"gitverse.ru/apavlov-systems/core-platform/internal/usecase"
	"gitverse.ru/apavlov-systems/core-platform/pkg/amqprpc"
)

func RegisterRoutes(server *amqprpc.Server, t usecase.Translation) {
	routes := &translationRoutes{t}

	// Регистрируем те же методы, что и в NATS
	server.Register("v1.translation.translate", routes.translate)
	server.Register("v1.translation.history", routes.history)
}
