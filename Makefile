DB_USER_NAME = "chats_service"
DB_NAME = "chats_service"


up:
	docker compose up --build -d --remove-orphans
down:
	docker compose down
logs:
	docker compose logs -f chats-service
dbshell:
	docker compose exec -it chats-service-db psql -U $(DB_USER_NAME) -d $(DB_NAME)
