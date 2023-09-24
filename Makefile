PROJECT = $(shell basename -s .git `git config --get remote.origin.url`)
REVISION ?= $(shell git rev-parse --short HEAD)
STACK_NAME ?= $(PROJECT)
BUCKET ?= $(USER)-$(STACK_NAME)
CF_TEMPLATE ?= deploy.yaml
PACKAGE_TEMPLATE = package.yaml

.PHONY: default test cover build validate tools
.PHONY: vars clean params bucket zip package
.PHONY: deploy destroy describe outputs

default: test

test:
	go test -v -race ./...

tidy:
	go mod tidy

cover:
	go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

vars:
	@echo PROJECT: $(PROJECT)
	@echo REVISION: $(REVISION)
	@echo STACK_NAME: $(STACK_NAME)
	@echo BUCKET: $(BUCKET)

.PHONY: cli
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

ENDPOINT ?= $(shell aws cloudformation describe-stacks \
	--stack-name $(STACK_NAME) \
	--query 'Stacks[].Outputs[?OutputKey==`ApiEndpoint`].OutputValue' \
	--output text)
.PHONY: remote-api-shorten
remote-api-shorten:
	curl -i $(ENDPOINT)/v1/shorten \
		-d '{"url":"https://go.dev"}'

.PHONY: remote-api-resolve
remote-api-resolve:
	curl -i $(ENDPOINT)/bn9Y9rho

## Deploy as an AWS Lambda function

bucket:
	-aws s3 mb s3://$(BUCKET)

clean:
	-@mkdir build
	-@rm -rf build/*
	-@rm package.yaml

build:
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -tags lambda.norpc -o ./build/shortener/bootstrap ./functions/shortener

zip:
	@cd ./build/shortener && zip bootstrap.zip bootstrap

tools:
	@which sam || brew install aws-sam-cli

package: build zip
	sam validate --template $(CF_TEMPLATE)
	sam package \
		--template-file $(CF_TEMPLATE) \
		--output-template-file $(PACKAGE_TEMPLATE) \
		--s3-bucket $(BUCKET)

deploy: clean package
	sam deploy \
		--template-file $(PACKAGE_TEMPLATE) \
		--stack-name $(STACK_NAME) \
		--capabilities CAPABILITY_IAM \
		--no-fail-on-empty-changeset \
		--disable-rollback

validate:
	aws cloudformation validate-template \
		--template-body file://$(CF_TEMPLATE)

destroy: clean
	-aws cloudformation delete-stack --stack-name $(STACK_NAME)

describe:
	aws cloudformation describe-stacks \
		--stack-name $(STACK_NAME)

outputs:
	aws cloudformation describe-stacks \
		--stack-name $(STACK_NAME) \
		--query 'Stacks[].Outputs'