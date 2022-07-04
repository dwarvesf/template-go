APP_NAME=api
DEFAULT_PORT=8100
.PHONY: setup init build dev test db-migrate-up db-migrate-down

setup:
	go install github.com/rubenv/sql-migrate/...@latest
	go install github.com/golang/mock/mockgen@v1.6.0
	go install github.com/vektra/mockery/v2@latest
	cp .env.sample .env
	make mock
	make init

init:
	make remove-infras
	docker-compose up -d
	@echo "Waiting for database connection..."
	@while ! docker exec db_local pg_isready -h localhost -p 5432 > /dev/null; do \
		sleep 1; \
	done
	make migrate-up
	make seed-db

remove-infras:
	docker-compose stop; docker-compose rm -f

build:
	env GOOS=darwin GOARCH=amd64 go build -o bin ./...

dev:
	go run ./cmd/server/main.go

test:
	make mock
	@PROJECT_PATH=$(shell pwd) go test -coverprofile=c.out -failfast -timeout 5m ./...


migrate-new:
	sql-migrate new -env=local ${name}

migrate-up:
	sql-migrate up -env=local

migrate-down:
	sql-migrate down -env=local

docker-build:
	docker build \
	--build-arg DEFAULT_PORT="${DEFAULT_PORT}" \
	-t ${APP_NAME}:latest .

seed-db:
	@docker cp data/seed/seed.sql  db_local:/seed.sql
	@docker exec -t db_local sh -c "PGPASSWORD=postgres psql -U postgres -d db_local -f /seed.sql"

reset-db:
	make init
	make migrate-up
	make seed-db

mock:
	mockery --dir pkg/entities --all --recursive --keeptree 
	mockery --dir pkg/repo --all --recursive --keeptree
