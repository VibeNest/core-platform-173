package amqprpc

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Handler func(ctx context.Context, data []byte) (interface{}, error)

type Server struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func (s *Server) Register(routingKey string, handler Handler) {
	// 1. Объявляем очередь (Queue), которую будем слушать
	q, err := s.channel.QueueDeclare(
		routingKey, // Имя очереди (например, "v1.translation.translate")
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		log.Printf("AMQP RPC: failed to declare a queue: %v", err)
		return
	}

	// 2. ИНИЦИАЛИЗИРУЕМ msgs (Создаем потребителя)
	// Именно здесь появляется переменная msgs, которой не хватало
	msgs, err := s.channel.Consume(
		q.Name, // имя очереди
		"",     // consumer
		true,   // auto-ack (автоматическое подтверждение получения)
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Printf("AMQP RPC: failed to register a consumer: %v", err)
		return
	}

	// 3. Запускаем обработку в фоновом режиме (горутине)
	go func() {
		for d := range msgs {
			// Вызываем хендлер контроллера
			res, err := handler(context.Background(), d.Body)
			if err != nil {
				log.Printf("AMQP RPC: handler error: %v", err)
				continue
			}

			// Кодируем ответ в JSON
			responseBytes, _ := json.Marshal(res)

			// Отправляем ответ обратно отправителю (Request-Reply)
			s.channel.Publish(
				"",        // exchange
				d.ReplyTo, // очередь для ответа (берем из запроса)
				false,
				false,
				amqp.Publishing{
					ContentType:   "application/json",
					CorrelationId: d.CorrelationId, // ID запроса для клиента
					Body:          responseBytes,
				},
			)
		}
	}()
}
