DOTF_VAR := programVersion
VERSION := $(shell git rev-parse --short HEAD)

.DEFAULT_GOAL := help

.PHONY: build-cli
build-cli: test ## Build the cli application
	go build -ldflags "-X main.$(DOTF_VAR)=$(VERSION)" -o bin/dotf-cli  cmd/dotf-cli/main.go

.PHONY: build-tray
build-tray: test ## Build the tray application
	go build -ldflags "-X main.$(DOTF_VAR)=$(VERSION)" -o bin/dotf-tray cmd/dotf-tray/main.go

.PHONY: build
build: build-cli build-tray ## Build all apps.

.PHONY: install-cli
install-cli: build-cli ## Installs cli app into default go location
	cd cmd/dotf-cli/ && go install -ldflags "-X main.$(DOTF_VAR)=$(VERSION)"

.PHONY: install-tray
install-tray: build-tray ## Installs tray app into default go location
	cd cmd/dotf-tray/ && go install -ldflags "-X main.$(DOTF_VAR)=$(VERSION)"

.PHONY: install-all
install-all: build-all install-cli install-tray ## Install all applications.

.PHONY: test
test: ## Run tests.
	go test -v ./pkg/...

.PHONY: install-ubuntu-deps
install-ubuntu-deps: ## Installs tray app deps for ubuntu
	sudo apt-get install gcc libgtk-3-dev libayatana-appindicator3-dev

.PHONY: help
help: ## Display this help
	@grep -E '^[a-zA-Z_0-9-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
