.PHONY: run-http
run-http:
	go run cmd/main.go http run

.PHONY: run-task
run-task:
	go run cmd/main.go task run

.PHONY: run-subscriber
run-subscriber:
	go run cmd/main.go subscriber run

.PHONY: go-build
go-build:
	go build -o ./.bin/app ./cmd

.PHONY: go-test
go-test:
	go test ./...

.PHONY: go-lint
go-lint:
	golangci-lint run -v --timeout=10m -c .golangci.yml

.PHONY: go-mod-download
go-mod-download:
	go mod download

.PHONY: go-mod-tidy
go-mod-tidy:
	go mod tidy

.PHONY: go-format
go-format:
	go fmt ./... && \
	gci write --skip-generated -s standard -s default -s "prefix(github.com/MoneyForest/go-clean-architecture-boilerplate)" .

.PHONY: swagger-gen
swagger-gen:
	go install github.com/swaggo/swag/cmd/swag@latest && \
	swag init -g cmd/main.go -o tools/swag && \
	swag fmt

.PHONY: sqlc-gen
sqlc-gen:
	sqlc generate

.PHONY: docker-clean-and-up
docker-clean-and-up:
	docker compose down -v && docker compose up -d --build

.PHONY: docker-clean
docker-clean:
	docker compose down -v

.PHONY: docker-up
docker-up:
	docker compose up -d --build

.PHONY: docker-down
docker-down:
	docker compose down

.PHONY: docker-logs
docker-logs:
	docker compose logs -f

.PHONY: clean
clean:
	rm -rf ./.bin
	rm -rf ./tools/swag
