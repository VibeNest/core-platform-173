Запуск проекта
```go
go test ./...
go build ./cmd/app
```

## Error Handling Strategy

Проект следует принципам **Clean Architecture** в обработке ошибок:

1. **Error Wrapping**: В слое инфраструктуры (Repository) ошибки оборачиваются через `%w`. Это сохраняет полную трассировку (стек) ошибки для логирования, не теряя исходный контекст БД.
   ```go
   // Пример из internal/repo/persistent
   return fmt.Errorf("HistoryRepo - StoreHistory - Exec: %w", err)

Observability vs Security: Технические детали (SQL errors, stack traces) логируются внутри системы для отладки, но никогда не экспонируются через API. Это гарантирует безопасность архитектуры при сохранении полной наблюдаемости для инженеров.

Посомтреть записи в БД
```BASH
docker exec -it core-platform-db psql -U user -d core_db -c "SELECT * FROM history;"
```

## 🛠 Проверка работоспособности (Testing Guide)

После запуска инфраструктуры (`docker-compose up -d`) и приложения (`make run`), вы можете проверить работу сервиса через разные протоколы. Все они используют единую бизнес-логику (**UseCase**) и базу данных.

### 1. REST API (HTTP / Fiber)
Самый простой способ проверки через терминал с помощью `curl`.

*   **Проверка статуса (Healthcheck):**
    ```bash
    curl http://localhost:8080/healthz
    ```
*   **Создать новый перевод:**
    ```bash
    curl -X POST http://localhost:8080/v1/translation/do \
      -H "Content-Type: application/json" \
      -d '{"source":"en", "destination":"ru", "original":"Hello World"}'
    ```
*   **Получить историю всех переводов:**
    ```bash
    curl http://localhost:8080/v1/translation/history
    ```

---

### 2. gRPC API
Для тестирования gRPC используйте утилиту [grpcurl](https://github.com). Сервер поддерживает **Reflection API**.

*   **Установка утилиты:**
   ```bash
      go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
   ```

*   **Список доступных методов:**
    ```bash
    grpcurl -plaintext localhost:8081 list v1.Translation
    ```
*   **Вызов метода Translate (JSON данные):**
    ```bash
    grpcurl -plaintext -d '{"source":"en", "destination":"ru", "original":"Hello gRPC"}' \
      localhost:8081 v1.Translation/Translate
    ```
*   **Запрос истории через gRPC:**
    ```bash
    grpcurl -plaintext localhost:8081 v1.Translation/GetHistory
    ```

### 3. NATS RPC (Message Bus)
Для проверки взаимодействия через NATS используется официальная утилита `nats`. В этом режиме сервис работает как RPC-сервер, принимая запросы из очереди и отправляя ответы во временные темы (Reply Subjects).
*   **Установка утилиты:**
   ```bash
      go install github.com/nats-io/natscli/nats@latest
   ```
*   **Выполнить перевод (Translate):**
    ```bash
    nats request v1.translation.translate '{"source":"en", "destination":"ru", "original":"Hello NATS"}'
    ```
    *Ожидаемый ответ: JSON с переведенным текстом и метаданными.*

*   **Запросить историю всех переводов (History):**
    ```bash
    nats request v1.translation.history '{}'
    ```
    *Ожидаемый ответ: JSON-массив со всеми записями из базы данных PostgreSQL.*

*   **Проверка Round Trip Time (RTT):**
    Обратите внимание на значение `Received with rtt` в выводе команды — оно показывает суммарное время на пересылку через брокер и обработку в Go-сервисе (обычно < 10ms).


*   **Примеры запросов и ответов:**
```bash
$  nats request v1.translation.translate '{"source":"en", "destination":"ru", "original":"Hello NATS"}'
    
14:05:56 Sending request on "v1.translation.translate"
14:05:56 Received with rtt 7.588638ms
{"id":"00000000-0000-0000-0000-000000000000","source":"en","destination":"ru","original":"Hello NATS","translation":"[EN-\u003eRU]: HELLO NATS","created_at":"0001-01-01T00:00:00Z"}

$   nats request v1.translation.history '{}'
14:06:05 Sending request on "v1.translation.history"
14:06:05 Received with rtt 2.674455ms
[{"id":"c39bfe7f-37f2-485d-ba58-616d5e9fd5a5","source":"en","destination":"ru","original":"Глубокое погружение в чистую архитектуру","translation":"[EN-\u003eRU]: ГЛУБОКОЕ ПОГРУЖЕНИЕ В ЧИСТУЮ АРХИТЕКТУРУ","created_at":"2026-03-17T19:39:13.034007+03:00"},{"id":"ec4751ee-2460-4d3f-aa11-99a9aed409ab","source":"en","destination":"ru","original":"Глубокое погружение в чистую архитектуру","translation":"[EN-\u003eRU]: ГЛУБОКОЕ ПОГРУЖЕНИЕ В ЧИСТУЮ АРХИТЕКТУРУ","created_at":"2026-03-17T19:31:05.167368+03:00"},{"id":"00000000-0000-0000-0000-000000000000","source":"en","destination":"ru","original":"","translation":"[EN-\u003eRU]: ","created_at":"0001-01-01T02:30:17+02:30"}]

```
