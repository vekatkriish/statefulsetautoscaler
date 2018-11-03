GOPATH=$(shell pwd)/vendor:$(shell pwd)
GOBIN=$(shell pwd)/bin
GOFILES=$(wildcard *.go)
GONAME=$(shell basename "$(PWD)")
PID=/tmp/go-$(GONAME).pid
APP_NAME=statefulset-operator
DOCKER_REPO=registry.rome.support
VERSION=$(shell ./version.sh)

all: watch

build:
	@echo "Building $(GOFILES) to ./bin"
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./bin/statefulset-operator cmd/statefulset-operator/main.go

get:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get .

install:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go install $(GOFILES)

watch: build stop start
	@fswatch -o *.go src/**/*.go | xargs -n1 -I{}  make restart || make stop

restart: stop clean build start

start:
	@echo "Starting bin/$(GONAME)"
	@./bin/$(GONAME) & echo $$! > $(PID)

stop:
	@echo "Stopping bin/$(GONAME)"
	@-kill `cat $(PID)` || true
clean:
	@echo "Cleaning"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean
build-docker: ## Build the container
	docker build -t $(APP_NAME):$(VERSION) .

build-docker-nc: ## Build the container without caching
	docker build --no-cache -t $(APP_NAME):$(VERSION) .
# Docker publish
publish: publish-version ## Publish the `{version}` ans `latest` tagged containers to ECR

publish-latest: tag-latest ## Publish the `latest` taged container to ECR
	@echo 'publish latest to $(DOCKER_REPO)'
	docker push $(DOCKER_REPO)/$(APP_NAME):latest

publish-version: tag-version ## Publish the `{version}` taged container to ECR
	@echo 'publish $(DOCKER_REPO)/$(APP_NAME):$(VERSION) to $(DOCKER_REPO)'
	docker push $(DOCKER_REPO)/$(APP_NAME):$(VERSION)

# Docker tagging
tag: tag-latest tag-version ## Generate container tags for the `{version}` ans `latest` tags

tag-latest: ## Generate container `{version}` tag
	@echo 'create tag latest'
	docker tag $(APP_NAME) $(DOCKER_REPO)/$(APP_NAME):latest

tag-version: ## Generate container `latest` tag
	@echo 'create tag $(VERSION)'
	docker tag $(APP_NAME) $(DOCKER_REPO)/$(APP_NAME):$(VERSION)
release: build-nc publish

.PHONY: build get install run watch start stop restart clean