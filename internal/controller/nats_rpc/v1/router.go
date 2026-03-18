package v1

import (
	"gitverse.ru/apavlov-systems/core-platform/internal/usecase"
	"gitverse.ru/apavlov-systems/core-platform/pkg/logger"
	"gitverse.ru/apavlov-systems/core-platform/pkg/natsrpc"
)

func RegisterRoutes(server *natsrpc.Server, t usecase.Translation, l *logger.Logger) {
	routes := newTranslationRoutes(t, l)

	server.Register("v1.translation.translate", routes.translate)
	server.Register("v1.translation.history", routes.history)
}
