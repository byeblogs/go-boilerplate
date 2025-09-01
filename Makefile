-include .env

APP_NAME := go-boilerplate
MIGRATION_DIR := ./migrations

DBMATE_URL := postgres://${POSTGRES_USERNAME}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB_NAME}?sslmode=disable

.PHONY: build
build:
	go build -v -o ./bin/${APP_NAME} ./internal/apps/rest

.PHONY: deploy
deploy:
	go build -v -o ./bin/${APP_NAME} ./internal/apps/rest
	    sudo supervisorctl restart ${APP_NAME}

.PHONY: start
start:
	chmod +x ./bin/${APP_NAME}
	./bin/${APP_NAME}

.PHONY: start-dev
start-dev:
	air -c .air.toml

.PHONY: run
run:
	go run ./internal/apps/rest

.PHONY: compile
compile:
	GOOS=linux GOARCH=386 go build -o ./bin/main-linux-386 ./internal/apps/rest
	GOOS=windows GOARCH=386 go build -o ./bin/main-windows-386 ./internal/apps/rest

.PHONY: migration-up
migration-up:
	migrate -database ${DBMATE_URL} -path migrations up

.PHONY: migration-down
migration-down:
	migrate -database ${DBMATE_URL} -path migrations down

.PHONY: migration-create
migration-create:
	@read -p "Enter the migration name: " MIGRATION_NAME; \
	migrate create -ext sql -dir $(MIGRATION_DIR) $$MIGRATION_NAME

.PHONY: migration-down-1
migration-down-1:
	migrate -database ${DBMATE_URL} -path migrations down 1

.PHONY: test
test:
	mkdir -p coverage
	go test -v -coverpkg=./internal/usecases/... -coverprofile ./coverage/coverage.out ./internal/usecases/... -count=1 | tee test_output.log
	go tool cover -html=./coverage/coverage.out -o ./coverage/coverage.html

.PHONY: generate-mock
generate-mock:
	if [ -d ./internal/mocks ]; then find ./internal/mocks -type f -delete; fi
	chmod +x ./scripts/generate_mocks.sh
	./scripts/generate_mocks.sh