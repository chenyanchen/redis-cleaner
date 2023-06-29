all: build

build:
	go build -o bin/redis-cleaner main.go
