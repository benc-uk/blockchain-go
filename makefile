# Common variables
VERSION := 0.0.1
BUILD_INFO := Manual build 

# Things you don't want to change
REPO_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
GOLINT_PATH := $(REPO_DIR)/bin/golangci-lint
AIR_PATH := $(REPO_DIR)/bin/air

.PHONY: help get-tools run-server build-server build-client lint lint-fix clean
.DEFAULT_GOAL := help

help: ## 💬 This help message :)
	@figlet $@ || true
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

get-tools: ## 🔮 Install dev tools into project bin directory
	@figlet $@ || true
	@$(GOLINT_PATH) > /dev/null 2>&1 || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin/
	@$(AIR_PATH) -v > /dev/null 2>&1 || curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh
	
lint: ## 🌟 Lint & format, will not fix but sets exit code on error
	@figlet $@ || true
	@$(GOLINT_PATH) > /dev/null || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh
	$(GOLINT_PATH) run ./...

lint-fix: ## 🔍 Lint & format, will try to fix errors and modify code
	@figlet $@ || true
	@$(GOLINT_PATH) > /dev/null || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh
	$(GOLINT_PATH) run ./... --fix

run-server: ## 🏃 Run server API locally, with hot reload
	@figlet $@ || true
	$(AIR_PATH) -c .air.toml

build-server: ## 🔨 Build server API, resulting binary in out/blockchain-api
	@figlet $@ || true
	go build -o out/blockchain-api api/main.go

build-client: ## 🔨 Build client command line tool, resulting binary in out/blockchain
	@figlet $@ || true
	go build -o out/blockchain client/main.go 

clean: ## 🧹 Clean up; database, binaries & temp files
	@figlet $@ || true
	@rm -rf ./tmp/*
	@rm -rf ./out/*
	@rm -rf ./*.db
