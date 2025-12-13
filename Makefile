
APP=todo

.PHONY: fmt lint test run tidy

fmt:
	go fmt ./...

tidy:
	go mod tidy

test:
	go test ./...

lint:
	golangci-lint run ./...

run:
	go run ./cmd/$(APP)
