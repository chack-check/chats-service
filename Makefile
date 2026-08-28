DB_USER_NAME = "chats_service"
DB_NAME = "chats_service"
GQL_DIR = "internal/infrastructure/api"


up:
	docker compose up --build -d --remove-orphans
down:
	docker compose down
api_logs:
	docker compose logs -f chats-service-api
grpcserver_logs:
	docker compose logs -f chats-service-grpcserver
consumer_logs:
	docker compose logs -f chats-service-consumer
dbshell:
	docker compose exec -it chats-service-db psql -U $(DB_USER_NAME) -d $(DB_NAME)
generate-gql:
	go run -C $(GQL_DIR) github.com/99designs/gqlgen generate
