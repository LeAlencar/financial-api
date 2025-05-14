.PHONY: migrate sqlc install-tools run

migrate:
	go run services/s2-processor/cmd/migrate/main.go

sqlc:
	sqlc generate

install-tools:
	go install github.com/jackc/tern/v2@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

run:
	go run services/s2-processor/cmd/api/main.go 