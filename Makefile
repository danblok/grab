build:
	@go build -o ./bin/grab ./cmd/main.go

test:
	go test -v -cover ./...

gobin = ~/go/bin
install:
	@go install ./cmd/main.go
	@mv $(gobin)/main $(gobin)/grab
	@tar -cf $(gobin)/grab.tar $(gobin)/grab >/dev/null 2>&1
	@tar -czf $(gobin)/grab.tar.gz $(gobin)/grab >/dev/null 2>&1
