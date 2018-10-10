.PHONY: all

BINARY_NAME=prometheus-p1-exporter

all: install build-binary build-binary-arm

install:
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	dep ensure

build-binary:
	go build -o dist/${BINARY_NAME} main.go

build-binary-arm:
	GOOS=linux GOARCH=arm go build -o dist/${BINARY_NAME}_arm main.go

test:
	go test -cover prometheus-p1-exporter/conn prometheus-p1-exporter/parser
