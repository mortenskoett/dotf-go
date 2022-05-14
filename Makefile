VNAME := programVersion
DOTF_VERSION := $(shell git rev-parse --short HEAD)

help:
	echo "Specify target"

build:
	echo hello $(VNAME)
	go build -ldflags "-X main.$(VNAME)=$(DOTF_VERSION)" -o bin/dotf-cli  cmd/dotf-cli/main.go 
	go build -ldflags "-X main.$(VNAME)=$(DOTF_VERSION)" -o bin/dotf-tray cmd/dotf-tray/main.go 

install: build
	cd cmd/dotf-cli/ && go install -ldflags "-X main.$(VNAME)=$(DOTF_VERSION)"
	cd cmd/dotf-tray/ && go install -ldflags "-X main.$(VNAME)=$(DOTF_VERSION)"

test: build
	go test -v ./pkg/...
