test:
	go test -v ./...

build:
	go build .

run:
	go run .

backend-run:
	fastapi dev tools/backend.py

.PHONY: test build run backend-run