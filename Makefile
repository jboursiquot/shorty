.PHONY: default test cover vendor

default: test

test:
	go test -v -race ./...

tidy:
	go mod tidy

cover:
	go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

cli:
	go run cmd/shortener/*.go -u https://go.dev