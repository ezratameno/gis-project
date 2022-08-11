SHELL := /bin/bash

build:
	go build -o gis-project ./cmd/

run: 
	go run ./cmd/web/

tidy:
	@go mod tidy
	@go mod verify
	@go mod vendor