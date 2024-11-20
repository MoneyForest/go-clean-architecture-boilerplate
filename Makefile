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
	gci write --skip-generated -s standard -s default -s "prefix(github.com/MoneyForest/go-clean-boilerplate)" .

.PHONY: swagger-gen
swagger-gen:
	swag init -g cmd/main.go -o pkg/swag && \
	swag fmt

.PHONY: docker-clean-and-up
docker-clean-and-up:
	docker-compose down -v && docker-compose up -d --build

.PHONY: docker-clean
docker-clean:
	docker-compose down -v

.PHONY: docker-up
docker-up:
	docker-compose up -d --build

.PHONY: docker-logs
docker-logs:
	docker-compose logs -f

.PHONY: clean
clean:
	rm -rf ./.bin
	rm -rf ./pkg/swag
