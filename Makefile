ifneq ($(wildcard .env),)
    include .env
    export $(shell sed 's/=.*//' .env)
endif

DOCKER_MIGRATE := docker run --rm --net=host -v $(shell pwd)/database/migrations:/migrations migrate/migrate

server:
	go run cmd/api/main.go

migrateup:
	$(DOCKER_MIGRATE) -path=/migrations -database "$(DATABASE_URL)" -verbose up

migratedown:
	$(DOCKER_MIGRATE) -path=/migrations -database "$(DATABASE_URL)" -verbose down -all

migrateforce:
	$(DOCKER_MIGRATE) -path=/migrations -database "$(DATABASE_URL)" force $(v)

create_migration:
	migrate create -ext sql -dir database/migrations -seq $(name)

sqlc:
	sqlc generate

upcompose:
	docker compose up -d

dropcompose:
	docker compose down -v

.PHONY: sqlc server upcompose dropcompose migrateup migratedown create_migration