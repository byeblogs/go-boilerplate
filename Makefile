.PHONY: clean test security swag build run \
        migrate.up migrate.down migrate.force migrate.create migrate.status \
        seed \
        docker.run docker.setup docker.fiber docker.fiber.build docker.postgres \
        docker.stop docker.stop.fiber docker.stop.postgres docker.dev docker.logs docker.reset

## ─── App Config ──────────────────────────────────────────────────────────────
APP_NAME = go-boilerplate
BUILD_DIR = ./build

APP_HOST = 0.0.0.0
APP_PORT = 9100
APP_READ_TIMEOUT = 30
APP_DEBUG = false

## ─── Database Config ─────────────────────────────────────────────────────────
DB_HOST = localhost
DB_PORT = 5432
DB_USER = killabyss
DB_PASS = password
DB_NAME = yaahtze_db
DB_SSL_MODE = disable

DATABASE_URL = postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)
MIGRATIONS_FOLDER = ./platform/migrations
MIGRATE = migrate -path $(MIGRATIONS_FOLDER) -database "$(DATABASE_URL)"

## ─── JWT Config ──────────────────────────────────────────────────────────────
JWT_SECRET_KEY = 97fb38c796601c55cf6f9dc57c3250f5
JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT = 15

## ─── Housekeeping ────────────────────────────────────────────────────────────
clean:
	rm -rf $(BUILD_DIR)/* cover.out *.out

security:
	gosec -quiet ./...

test: security
	go test -v -timeout 30s -coverprofile=cover.out -cover ./...
	go tool cover -func=cover.out

swag:
	swag init

build: swag clean
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) main.go

run: build
	$(BUILD_DIR)/$(APP_NAME)

## ─── Database: Migrations & Seeds ────────────────────────────────────────────
migrate.up:
	$(MIGRATE) up

migrate.down:
	$(MIGRATE) down

migrate.force:
	@if [ -z "$(version)" ]; then \
		echo "❌ Please provide a version: make migrate.force version=<N>"; \
		exit 1; \
	fi
	$(MIGRATE) force $(version)

migrate.create:
	@if [ -z "$(name)" ]; then \
		echo "❌ Please provide a name: make migrate.create name=add_users"; \
		exit 1; \
	fi
	migrate create -ext sql -dir $(MIGRATIONS_FOLDER) -seq $(name)

migrate.status:
	$(MIGRATE) version

seed:
	PGPASSWORD=$(DB_PASS) psql -h $(DB_HOST) -p $(DB_PORT) -U$(DB_USER) -d $(DB_NAME) -a -f platform/seeds/001_seed_user_table.sql
	PGPASSWORD=$(DB_PASS) psql -h $(DB_HOST) -p $(DB_PORT) -U$(DB_USER) -d $(DB_NAME) -a -f platform/seeds/002_seed_book_table.sql

## ─── Dockerized Dev Env ──────────────────────────────────────────────────────
docker.run: docker.setup docker.postgres docker.fiber migrate.up
	@echo "\n===========FGB==========="
	@echo "App is running... Visit: http://localhost:5000 OR http://localhost:5000/swagger/"

docker.setup:
	docker network inspect dev-network >/dev/null 2>&1 || \
		docker network create -d bridge dev-network
	docker volume create fibergb-pgdata

docker.fiber.build: swag
	docker build -t fibergb:latest .

docker.fiber: docker.fiber.build
	docker run --rm -d \
		--name fibergb-api \
		--network dev-network \
		-p 5000:5000 \
		fibergb

docker.postgres:
	docker run --rm -d \
		--name fibergb-postgres \
		--network dev-network \
		-e POSTGRES_USER=$(DB_USER) \
		-e POSTGRES_PASSWORD=$(DB_PASS) \
		-e POSTGRES_DB=$(DB_NAME) \
		-v fibergb-pgdata:/var/lib/postgresql/data \
		-p $(DB_PORT):5432 \
		postgres

docker.stop: docker.stop.fiber docker.stop.postgres

docker.stop.fiber:
	docker stop fibergb-api || true

docker.stop.postgres:
	docker stop fibergb-postgres || true

docker.dev:
	docker-compose up

docker.logs:
	docker logs -f fibergb-api

docker.reset: docker.stop
	docker volume rm -f fibergb-pgdata || true
