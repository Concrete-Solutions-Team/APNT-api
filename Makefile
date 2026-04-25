app-up:
	docker compose up -d --build

app-down:
	docker compose down -v

app-reload: app-down app-up

NAME ?= Alice

api-register:
	curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username": "$(NAME)", "password": "password1", "role": "student"}' \
  -c cookies.txt

api-login:
	curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "$(NAME)", "password": "password1", "role": "student"}' \
  -c cookies.txt

api-me:
	curl -X GET http://localhost:8080/auth/me \
  -b cookies.txt