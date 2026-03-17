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