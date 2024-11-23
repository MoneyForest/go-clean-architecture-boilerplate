include misc/make/docker.make
include misc/make/go.make
include misc/make/migrate.make

.PHONY: run-http
run-http:
	go run cmd/main.go http run

.PHONY: run-task
run-task:
	go run cmd/main.go task run

.PHONY: run-subscriber
run-subscriber:
	go run cmd/main.go subscriber run

.PHONY: gen-swagger
gen-swagger:
	go install github.com/swaggo/swag/cmd/swag@latest && \
	swag init -g cmd/main.go -o tools/swag && \
	swag fmt

.PHONY: gen-sqlc
gen-sqlc:
	sqlc generate

.PHONY: clean
clean:
	rm -rf ./.bin
	rm -rf ./tools/swag

.PHONY: setup
setup: docker-setup go-setup db-migrate-setup clean gen-swagger gen-sqlc
