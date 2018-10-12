.PHONY: install install-deps build-binary build-binary-arm install-binary test

BINARY_NAME=prometheus-p1-exporter

install: install-deps build-binary install-binary

install-deps:
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	dep ensure

build-binary:
	go build -o dist/${BINARY_NAME} main.go

build-binary-arm:
	GOOS=linux GOARCH=arm go build -o dist/${BINARY_NAME}_arm main.go

install-binary:
	mv dist/${BINARY_NAME} /usr/local/bin/${BINARY_NAME}
	chmod +x /usr/local/bin/${BINARY_NAME}
	@echo "Installed binary at '/usr/local/bin/${BINARY_NAME}'"

test:
	go test -cover -coverprofile=cover.out github.com/jordyv/prometheus-p1-exporter/conn github.com/jordyv/prometheus-p1-exporter/parser
