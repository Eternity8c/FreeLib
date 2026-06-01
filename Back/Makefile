include .env
export

export PROJECT_ROOT=$(shell pwd)

# Запуск контейнера
env-up:
	@docker-compose up freelib-postgres -d

# Остановка контейнера
env-down:
	@docker-compose down freelib-postgres

# Очистка окружения
env-cleanup:
	@read -p "Очистить все volume файлы окружения? [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker-compose down && \
		rm -rf out/pgdata && \
		echo "Файлы окружения очищены"; \
	else \
		echo "Очистка окружения отменена"; \
	fi

# Запуск приложения
run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run cmd/server/main.go

# Запуск тестов
test: 
	go test ./... -v

# Создание миграций
migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Отсутсвует необходимый параметр 'seq'. Пример: make migrate-create seq=init"; \
		exit 1; \
	fi; \

	docker-compose run --rm freelib-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

# Применение миграций 
migrate-up:
	@make migrate-action action=up

# Откат миграций
migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Отсутсвует необходимый параметр 'action'. Пример: migrate-action action=up"; \
		exit 1; \
	fi; \

	@docker-compose run --rm freelib-postgres-migrate \
		-path ./migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@freelib-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		$(action)