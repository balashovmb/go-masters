Запуск
docker compose -f docker-compose.yml up
go run main.go

Получить оценки
Фильтр по пользователю
curl -X GET "http://localhost:8080/reviews" -H "Accept: application/json" -d '{"filter": "user", "id": 1}'
Фильтр по объекту
curl -X GET "http://localhost:8080/reviews" -H "Accept: application/json" -d '{"filter": "object", "id": 1}'

Добавить оценку
curl -X POST http://localhost:8080/reviews -H "Content-Type: application/json" -d '{"user_id": 1, "object_id": 1, "text": "cool"}'

Средняя оценка
curl -X GET "http://localhost:8080/reviews/1/average" -H "Accept: application/json"
