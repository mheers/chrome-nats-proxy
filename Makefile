## This is a self-documented Makefile. For usage information, run `make help`:
##
## For more information, refer to https://suva.sh/posts/well-documented-makefiles/

SHELL := /bin/bash

include .env

all: help

##@ Running
start-server: ## Start the server
	@echo "Starting server..."
	docker-compose up nats

create-seeds:
	@echo "Creating seeds..."
	docker run --rm mheers/nats-seeder seeds

create-credentials:
	docker run --rm mheers/nats-seeder user-nkey \
	--operator-seed $(OPERATOR_SEED) \
	--account-seed $(ACCOUNT_SEED) \
	-u chrome-nats-proxy \
	-p "\$JS.API.>" -s "\$JS.API.>" -p "_INBOX.>" -s "_INBOX.>" \
	-p "instance.clipboard.*" \
	-s "instance.clipboard.*" \
	-p "instance.chrome-nats-proxy.*" \
	-s "instance.chrome-nats-proxy.*" \
	> mq.creds

build: ## Build the binary
	go build -o chrome-nats-proxy.bin .

build-windows:
	GOOS=windows GOARCH=amd64 go build -o chrome-nats-proxy.exe .

##@ Helpers

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[0-9a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
