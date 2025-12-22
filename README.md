# Go Microservice

Простой демонстрационный микросервис на Go 1.22 с CRUD по пользователям, rate limiting, метриками Prometheus и контейнеризацией.

## Запуск локально
```bash
go mod tidy
go run .
```
Сервис поднимется на `http://localhost:8080`.

## Эндпоинты
- `GET /api/users` — список пользователей
- `GET /api/users/{id}` — получить пользователя
- `POST /api/users` — создать пользователя (`{"name":"John","email":"john@example.com"}`)
- `PUT /api/users/{id}` — обновить пользователя
- `DELETE /api/users/{id}` — удалить
- `GET /metrics` — метрики Prometheus

## Rate limiting
1000 rps с burst 5000 (`golang.org/x/time/rate`). Ошибка 429 при превышении.

## Контейнеризация
```bash
docker build -t go-microservice .
docker run -p 8080:8080 go-microservice
```

### docker-compose + MinIO
```bash
docker-compose up --build
```
MinIO доступен на `:9000` (консоль `:9001`). Параметры подаются через переменные окружения.

## Нагрузочное тестирование
Пример команды:
```bash
wrk -t12 -c500 -d60s http://localhost:8080/api/users
```

## Дополнительно
- Асинхронные аудит-логи и уведомления отправляются в goroutine.
- Грейсфул shutdown с таймаутом 5s.


