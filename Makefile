GOPROXY=https://proxy.golang.org,direct
DOCKER_COMPOSE_FILENAME=deployments/docker-compose/docker-compose.yaml

.DEFAULT_GOAL := _help

# Help: Each target starting with "_" will be ignored. Add description for target by adding "## <my help description>"
_help:
	@grep -E '^[0-9a-zA-Z][0-9a-zA-Z_-]+:.*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*## ?"}; {split($$1, t, ":"); printf "\033[1;34m%-30s\033[0m\t%s\n", t[1], $$2}'

# Running locally
start: _dot-env vendor-install ## Start the application and all dependencies
	 docker-compose --file ${DOCKER_COMPOSE_FILENAME} --profile httpapi up -d

start-from-scratch: _dot-env vendor-remove vendor-install ## Start the application and all dependencies with complete rebuild
	docker-compose --file ${DOCKER_COMPOSE_FILENAME} --profile httpapi up --remove-orphans --renew-anon-volumes --force-recreate --build -d

stop: _dot-env vendor-install ## Stop the application and all dependencies
	 docker-compose --file ${DOCKER_COMPOSE_FILENAME} --profile httpapi stop

logs-httpapi: ## Displays log output from HTTP API running locally
	docker-compose --file ${DOCKER_COMPOSE_FILENAME} --profile httpapi logs -f httpapi

start-swagger: ## Start swagger
	docker-compose --file ${DOCKER_COMPOSE_FILENAME} --profile swagger up -d

# Development
generate-go: ## Generates mocks and other stuff
	docker run --rm -v ${PWD}:/app -w /app docker.io/code7unner/alpine3.20-golang-1.21 \
	-c 'go generate ./...'

# Tests
test-unit: vendor-install ## Run unit tests with coverage
	docker run -v $(shell pwd):/app -w /app --rm docker.io/code7unner/alpine3.20-golang-1.21 \
	-c 'go test -count=1 -v -covermode count -coverpkg=./... -coverprofile=coverage.txt $$(go list ./... | grep -v /test/) && go tool cover -html=coverage.txt -o=coverage.html && go tool cover -func coverage.txt | sed -E "s/total:[^0-9]+/coverage-total: /" && cat coverage.txt | grep -E "mode:|internal(/domain/|/entity/|/usecase/)" | go tool cover -func=/dev/stdin | sed -E "s/total:[^0-9]+/coverage-entities-usecases: /" | grep coverage-entities-usecases'

# Linters
lint: lint-golangci lint-open-api-spec lint-architecture ## Start all linters

lint-golangci: vendor-install ## Start lint-golangci linter
	docker run --rm -v ${PWD}:/app -w /app golangci/golangci-lint:v1.53.1 \
    golangci-lint run --timeout 5m30s -v

lint-open-api-spec: ## Start lint-open-api-spec linter
	docker run -v ${PWD}:/app -w /app --rm jeanberu/swagger-cli:4.0.4 \
	swagger-cli validate /app/api/openapi-spec/httpapi.openapi.yaml

lint-architecture: ## Run Architecture checks for the project: dependencies, package content, cyclic dependencies, functions properties, naming rules [arch-go]
	docker run --rm -v $(shell pwd):/app -w /app docker.io/code7unner/alpine3.20-golang-1.21 \
	-c "arch-go -v"

vendor-reinstall: vendor-remove vendor-install ## Remove vendors and install from scratch

vendor-install: ## Install vendors
	@if [ -d "vendor" ]; then \
	    echo "Vendor folder already exists. Skip vendor installing."; \
	else \
	  	echo "Vendor installing...."; \
		docker run --rm -v ${PWD}:/app -v ${GOPATH}/pkg/mod:/go/pkg/mod -e GOPROXY=${GOPROXY} -w /app \
			 docker.io/code7unner/alpine3.20-golang-1.21 -c "go mod tidy && go mod vendor" && \
		echo "Vendor installing has been finished" ; \
	fi

vendor-remove: ## Remove vendors
	rm -rf vendor || true

_dot-env:
	touch deployments/docker-compose/.env
