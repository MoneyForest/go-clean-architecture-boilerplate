DB_USER := root
DB_PASSWORD := password
DB_HOST := localhost
DB_PORT := 3306
DB_NAME := maindb
TEST_DB_NAME := testdb
MYSQL_DSN := mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)?multiStatements=true
TEST_MYSQL_DSN := mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(TEST_DB_NAME)?multiStatements=true
MIGRATIONS_PATH := ./internal/infrastructure/gateway/mysql/schema

.PHONY: wait-for-mysql
wait-for-mysql:
	@echo "Waiting for MySQL to be ready..."
	@for i in $$(seq 1 30); do \
		if docker compose exec mysql mysqladmin ping -h localhost -u$(DB_USER) -p$(DB_PASSWORD) >/dev/null 2>&1; then \
			echo "MySQL is ready!"; \
			break; \
		fi; \
		echo "Waiting for MySQL... $$i"; \
		sleep 1; \
	done

.PHONY: db-migrate-install
db-migrate-install:
	go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

.PHONY: db-migrate-up
db-migrate-up: wait-for-mysql db-migrate-install
	migrate -database '$(MYSQL_DSN)' -path $(MIGRATIONS_PATH) up

.PHONY: test-db-migrate-up
test-db-migrate-up: wait-for-mysql db-migrate-install
	migrate -database '$(TEST_MYSQL_DSN)' -path $(MIGRATIONS_PATH) up

.PHONY: db-migrate-down
db-migrate-down: wait-for-mysql db-migrate-install
	migrate -database '$(MYSQL_DSN)' -path $(MIGRATIONS_PATH) down

.PHONY: db-migrate-version
db-migrate-version: wait-for-mysql db-migrate-install
	migrate -database '$(MYSQL_DSN)' -path $(MIGRATIONS_PATH) version

.PHONY: db-migrate-create
db-migrate-create: wait-for-mysql db-migrate-install
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $$name

.PHONY: db-migrate-setup
db-migrate-setup: wait-for-mysql db-migrate-install db-migrate-up

.PHONY: test-db-migrate-setup
test-db-migrate-setup: wait-for-mysql db-migrate-install test-db-migrate-up
