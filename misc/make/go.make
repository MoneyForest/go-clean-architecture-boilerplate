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

.PHONY: go-import-violation-check
go-import-violation-check:
	go run tools/import-violation-checker/main.go

.PHONY: go-setup
go-setup: go-mod-download go-mod-tidy
