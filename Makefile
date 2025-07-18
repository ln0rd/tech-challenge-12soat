BINARY_NAME=tech-challenge-12soat
MAIN_PATH=cmd

.PHONY: all build run clean test lint setup down run-bin

build:
	go build -o $(BINARY_NAME) $(MAIN_PATH)/main.go

run:
	go run $(MAIN_PATH)/main.go

test:
	go test ./...

lint:
	golangci-lint run

clean:
	go clean
	rm -f $(BINARY_NAME)

setup:
	docker compose up -d

down:
	docker compose down

run-bin: build
	./$(BINARY_NAME)
