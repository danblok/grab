build:
	@go build -o ./bin/grab ./cmd/main.go

test:
	go test -v -cover ./...
