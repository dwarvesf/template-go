.PHONY: build dev

build:
	go build -o server cmd/server/*.go

dev: build
	./server; rm server
