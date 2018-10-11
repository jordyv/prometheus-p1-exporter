.PHONY: all

BINARY_NAME=prometheus-p1-exporter

all: build-binary build-binary-arm

build-binary:
	go build -o dist/${BINARY_NAME} main.go

build-binary-arm:
	GOOS=linux GOARCH=arm go build -o dist/${BINARY_NAME}_arm main.go

test:
	go test -cover prometheus-p1-exporter/conn prometheus-p1-exporter/parser
