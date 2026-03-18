package main

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// 1. Подключаемся к RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	// 2. Создаем временную очередь для получения ответа (exclusive)
	msgs, err := ch.Consume(
		"amq.rabbitmq.reply-to", // Специальная очередь для Direct Reply-To
		"", true, true, false, false, nil,
	)
	if err != nil {
		log.Fatalf("failed to register a consumer: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 3. Публикуем запрос
	body := `{"source":"en", "destination":"ru", "original":"Hello Rabbit RPC"}`
	err = ch.PublishWithContext(ctx,
		"",                         // exchange
		"v1.translation.translate", // routing key (имя очереди в твоем сервисе)
		false, false,
		amqp.Publishing{
			ContentType: "application/json",
			ReplyTo:     "amq.rabbitmq.reply-to", // Указываем, куда слать ответ
			Body:        []byte(body),
		})
	if err != nil {
		log.Fatalf("failed to publish a message: %v", err)
	}

	fmt.Printf(" [x] Sent request: %s\n", body)

	// 4. Ждем ответ
	select {
	case d := <-msgs:
		fmt.Printf(" [.] Got response: %s\n", d.Body)
	case <-ctx.Done():
		log.Fatalf("timeout: no response from server")
	}
}
