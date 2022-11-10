POSTGRES_DB ?= yellow-jersey
POSTGRES_URL ?= postgres://postgres:postgres@0.0.0.0:5432/$(POSTGRES_DB)?sslmode=disable

.PHONY : init migrate-up clean-env

down:
	@docker-compose down
	@docker volume rm --force yellow-jersey_mongodb_data_container

init:
	@docker-compose up -d
	@sleep 3
	@docker-compose up -d --remove-orphans


migrate-up:
	docker run -v "$(PWD)/migrations/migrations:/migrations" --network host migrate/migrate -path=/migrations -database $(POSTGRES_URL) up $(TIMES)