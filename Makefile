ifneq ($(wildcard .env),)
    include .env
    export $(shell sed 's/=.*//' .env)
endif

DOCKER_MIGRATE := docker run --rm --net=host -v $(shell pwd)/database/migrations:/migrations migrate/migrate
DOCKER_SQLC := docker run --rm -v $(shell pwd):/src -w /src sqlc/sqlc

# --------------------------------------------------------------------------------------------------------------------
setup:
	make composeup && make migrateup && make sqlc

dev:
	go run cmd/api/main.go

server:
	GIN_MODE=release go run cmd/api/main.go

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./bin/bs-aesthetics-api ./cmd/api
# --------------------------------------------------------------------------------------------------------------------

migrateup:
	$(DOCKER_MIGRATE) -path=/migrations -database "$(DATABASE_URL)" -verbose up

migratedown:
	$(DOCKER_MIGRATE) -path=/migrations -database "$(DATABASE_URL)" -verbose down -all

migrateforce:
	$(DOCKER_MIGRATE) -path=/migrations -database "$(DATABASE_URL)" force $(v)

create_migration:
	$(DOCKER_MIGRATE) create -ext sql -dir /migrations -seq $(name)

sqlc:
	$(DOCKER_SQLC) generate

composeup:
	docker compose up -d

composedown:
	docker compose down

dropcompose:
	docker compose down -v

.PHONY: sqlc dev server composeup composedown dropcompose migrateup migratedown create_migration setup build