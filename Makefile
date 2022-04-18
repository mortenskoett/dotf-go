help:
	echo "Specify target"

build:
	go build -o bin/dotf-cli cmd/dotf-cli/main.go 
	go build -o bin/dotf-tray cmd/dotf-tray/main.go 

install: build
	cd cmd/dotf-cli/ && go install
	cd cmd/dotf-tray/ && go install

test: build
	go test -v ./pkg/...
