.PHONY: build
build:
	go build -v ./cmd/app_email

.DEFAULT_GOAL := build