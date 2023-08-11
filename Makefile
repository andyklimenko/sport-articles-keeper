test: lint
	@go test -count=1 -race ./...

build: lint
	@go build -o articles-keeper main.go

run-server: lint
	@go run main.go

lint: ## Lint the files
	@golangci-lint run
