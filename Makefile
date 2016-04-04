.PHONY: build, cmd

cmd:
	go build -o build/relay ./cmd/relay

build:
	docker build -t lunohq/relay .
