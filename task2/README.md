# Сервис, предоставляющий информацию о курсах валют
## Использование 
1. Запустить сервер

`go run cmd/main.go -p PORT_NUMBER`

PORT_NUMBER - номер порта. Порт по умолчанию 8000

2. Отправить запрос 

`curl -X POST http://localhost:8000/rate \
-H "Content-Type: application/json" \
-d '{"from": "CURRENCY1", "to": "CURRENCY2"}'`

CURRENCY1, CURRENCY2 - коды валют. Поддерживаются: USD, EUR, RUB
