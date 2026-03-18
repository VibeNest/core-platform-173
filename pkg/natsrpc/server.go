package natsrpc

import (
	"context"
	"encoding/json"
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
		// 1. Логируем входящий запрос (для отладки)
		log.Printf("NATS RPC: received request on [%s]", subject)

		// 2. Вызываем бизнес-логику
		res, err := handler(context.Background(), msg.Data)
		if err != nil {
			log.Printf("NATS RPC Error on %s: %v", subject, err)
			// Отправляем ошибку клиенту, чтобы он не висел по таймауту
			msg.Respond([]byte(`{"error": "internal error"}`))
			return
		}

		// 3. Превращаем результат (interface{}) в байты JSON
		responseBytes, err := json.Marshal(res)
		if err != nil {
			log.Printf("NATS RPC Marshal Error: %v", err)
			return
		}

		// 4. САМОЕ ВАЖНОЕ: Отправляем ответ обратно клиенту
		err = msg.Respond(responseBytes)
		if err != nil {
			log.Printf("NATS RPC Respond Error: %v", err)
		}
	})

	if err != nil {
		log.Printf("NATS RPC: failed to subscribe to %s: %v", subject, err)
	}
}
