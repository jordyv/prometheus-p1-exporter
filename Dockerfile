FROM golang:1.18-alpine AS builder
WORKDIR $GOPATH/src/app/
RUN apk add --no-cache git
ADD . .
RUN go get
RUN go build -o dist/prometheus-p1-exporter main.go

FROM alpine
COPY --from=builder /go/src/app/dist/prometheus-p1-exporter /usr/local/bin/prometheus-p1-exporter
RUN chmod +x /usr/local/bin/prometheus-p1-exporter
ENTRYPOINT ["prometheus-p1-exporter", "-listen", "0.0.0.0:8888"]
