package natsrpc

import (
	"context"
	"log"

	"github.com/nats-io/nats.go"
)

// Handler — это сигнатура функции, которую мы будем писать в контроллерах
// Она принимает сырые байты (JSON/Proto) и возвращает ответ или ошибку
type Handler func(ctx context.Context, data []byte) (interface{}, error)

// Server — структура нашего NATS RPC движка
type Server struct {
	conn *nats.Conn
}

// NewServer создает новый экземпляр сервера
func NewServer(conn *nats.Conn) *Server {
	return &Server{conn: conn}
}

// Register связывает тему (subject) в NATS с конкретной функцией-обработчиком
func (s *Server) Register(subject string, handler Handler) {
	_, err := s.conn.Subscribe(subject, func(msg *nats.Msg) {
		// Вызываем бизнес-логику через handler
		res, err := handler(context.Background(), msg.Data)
		if err != nil {
			log.Printf("NATS RPC Error on %s: %v", subject, err)
			// Здесь можно отправить структурированную ошибку в ответ
			return
		}

		// Пока просто отправляем ответ обратно (в будущем добавим JSON маршалинг)
		// s.conn.Publish(msg.Reply, responseBytes)
		_ = res
	})

	if err != nil {
		log.Printf("NATS RPC: failed to subscribe to %s: %v", subject, err)
	}
}
