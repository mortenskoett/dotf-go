all:

build:
	rm bin/*
	go build -o bin/dotf cmd/dotf-cli/main.go 
	go build -o bin/dotf-tray cmd/dotf-tray/main.go 
