all:

build:
	go build -o bin/dotf-cli cmd/dotf-cli/main.go 
	go build -o bin/dotf-tray cmd/dotf-tray/main.go 
