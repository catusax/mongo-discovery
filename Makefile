GOPATH:=$(shell go env GOPATH)

.PHONY: update
update:
	@go get -u

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: build
build:
	@go build -o mongo-discovery *.go

.PHONY: test
test:
	@go test -v ./... -cover

.PHONY: docker
docker:
	@docker build -t ghcr.io/catusax/mongo-discovery:latest .

.PHONY: docker-push
docker-push:
	@docker push ghcr.io/catusax/mongo-discovery:latest

