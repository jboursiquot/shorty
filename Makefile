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

.PHONY: run-local-api
run-local-api:
	go run functions/shortener/*.go -run-as-lambda=false

.PHONY: local-api-shorten
local-api-shorten:
	curl -i http://localhost:8080/v1/shorten \
		-d '{"url":"https://go.dev"}'

.PHONY: local-api-resolve
local-api-resolve:
	curl -i http://localhost:8080/bn9Y9rho

.PHONY: api-call
ENDPOINT ?= https://$(STACK_NAME).execute-api.us-east-1.amazonaws.com/prod
api-call:
	curl -i $(ENDPOINT)