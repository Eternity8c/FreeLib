include .env
export

export PROJECT_ROOT=$(shell pwd)

# Запуск докера
start-docker:
	@sudo systemctl start docker

# Запуск контейнера
env-up:
	@docker compose up freelib-postgres -d

# Остановка контейнера
env-down:
	@docker compose down freelib-postgres

# Очистка окружения
env-cleanup:
	@read -p "Очистить все volume файлы окружения? [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down && \
		rm -rf out/pgdata && \
		echo "Файлы окружения очищены"; \
	else \
		echo "Очистка окружения отменена"; \
	fi

# Запуск приложения
run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	go run cmd/freelib/main.go

# Запуск тестов
test: 
	go test ./... -v

# Создание миграций
migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Отсутсвует необходимый параметр 'seq'. Пример: make migrate-create seq=init"; \
		exit 1; \
	fi; \

	docker compose run --rm freelib-postgres-migrate \
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

	@dockerтcompose run --rm freelib-postgres-migrate \
		-path ./migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@freelib-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		$(action)

freelib-deploy:
	@docker compose up -d --build freelib

ps:
	@docker compose ps

swagger-gen:
	@docker compose run --rm swagger \
		init \
		-g cmd/freelib/main.go \
		-o docs \
		--parseInternal \
		--parseDependency
	@sudo chown -R $$(id -u):$$(id -g) docs
