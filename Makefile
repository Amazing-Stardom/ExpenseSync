.PHONY: run test cover swagger build docker-build docker-run docker-up docker-down

run:
	@if command -v air > /dev/null 2>&1; then \
		air; \
	else \
		go run main.go; \
	fi

test:
	go test ./... -v

cover:
	go test ./... -cover

swagger:
	swag init --output ./api-docs --packageName docs

build:
	go build -o expensesync main.go

docker-build:
	docker build -t ExpenseSync:latest .

docker-run:
	docker run -p 8080:8080 ExpenseSync:latest

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down
