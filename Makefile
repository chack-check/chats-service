dev:
	docker compose -f docker-compose.dev.yml down
	docker compose -f docker-compose.dev.yml up --build
test:
	docker compose -f docker-compose.test.yml up --build -d
	docker compose -f docker-compose.test.yml exec -it test-service go test -v ./... && docker compose -f docker-compose.test.yml down -v
	docker compose -f docker-compose.test.yml down -v
