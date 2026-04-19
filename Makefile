include .env
export

export PROJECT_ROOT=$(shell pwd)


swag:
	swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal

env-up:
	@docker compose up -d happytail-postgres

env-down:
	@docker compose down happytail-postgres

env-cleanup:
	@read -p "Очистить все volume файлы окружения ? Потеря данных. (y/n) " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down down happytail-postgres && \
		rm -rf out/pgdata && \
		echo "Database volume removed."; \
	else \
		echo "Очистка окружения отменена."; \
	fi



migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Отсутствует аргумент seq. Использование: make migrate-create seq=название_миграции"; \
		exit 1; \
	fi; \
	docker compose run --rm happytail-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@make migrate-action action=up	

migrate-down:
	make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Отсутствует аргумент action. Использование: make migrate-action action=up|down"; \
		exit 1; \
	fi; \
	docker compose run --rm happytail-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@happytail-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"

dev:
	docker compose up --build happytail-api