POSTGRES_DB ?= yellow-jersey
POSTGRES_URL ?= postgres://postgres:postgres@0.0.0.0:5432/$(POSTGRES_DB)?sslmode=disable
TIMES ?=

version:
	migrate -database $(POSTGRES_URL) -path ./migrations version

migrate-up:
	@migrate -database $(POSTGRES_URL) -path ./migrations up $(TIMES)

migrate-down:
	@migrate -database $(POSTGRES_URL) -path ./migrations down $(TIMES)

docker-exec-postgres:
	docker exec -it mirror_postgres psql -U postgres -d $(POSTGRES_DB)

install:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

test:
	go test ./...
