BINARY_NAME=tech-challenge-12soat
MAIN_PATH=cmd

.PHONY: all build run clean test lint setup up down run-bin sonar-up sonar-down sonar-logs sonar-init-token sonar-scan

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
	docker compose up -d postgres sonarqube-db sonarqube

up: setup

down:
	docker compose down

run-bin: build
	./$(BINARY_NAME)

# SonarQube helpers
sonar-up:
	docker compose up -d sonarqube-db sonarqube

sonar-down:
	docker compose rm -sfv sonarqube sonarqube-db || true

sonar-logs:
	docker compose logs -f sonarqube | cat

# First-run helper to print the default credentials and URL
sonar-init-token:
	@echo "Open http://localhost:9000 (user: admin / password: admin)."
	@echo "Create a token in your profile and export SONAR_TOKEN=your_token"

# Run local analysis with sonar-scanner (requires SONAR_TOKEN and sonar-scanner installed)
sonar-scan:
	@if [ -z "$$SONAR_TOKEN" ]; then echo "SONAR_TOKEN not set. Export it first."; exit 1; fi
	sonar-scanner -Dsonar.login=$$SONAR_TOKEN -Dsonar.host.url=http://localhost:9000
