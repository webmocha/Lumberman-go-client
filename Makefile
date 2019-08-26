build:
	go build -o bin/lmc .

build-linux64:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/lmc .

run:
	./bin/lmc

dev: build run

.PHONY: build build-linux64 run dev
