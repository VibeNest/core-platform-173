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
