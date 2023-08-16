.PHONY : up down

down:
	@docker-compose down
	@docker volume rm --force yellow-jersey_mongodb_data_container

up:
	@docker-compose up -d --remove-orphans