DOTF_VAR := programVersion
VERSION := $(shell git rev-parse --short HEAD)

help:
	echo "Specify target"

build:
	echo hello $(DOTF_VAR)
	go build -ldflags "-X main.$(DOTF_VAR)=$(VERSION)" -o bin/dotf-cli  cmd/dotf-cli/main.go
	go build -ldflags "-X main.$(DOTF_VAR)=$(VERSION)" -o bin/dotf-tray cmd/dotf-tray/main.go

install: build
	cd cmd/dotf-cli/ && go install -ldflags "-X main.$(DOTF_VAR)=$(VERSION)"
	cd cmd/dotf-tray/ && go install -ldflags "-X main.$(DOTF_VAR)=$(VERSION)"

test: build
	go test -v ./pkg/...
