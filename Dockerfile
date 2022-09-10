FROM golang:1.19.0 as builder
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download

# Build
COPY . .
RUN git rev-parse --short HEAD
RUN GIT_COMMIT=$(git rev-parse --short HEAD) && \
    CGO_ENABLED=0 go build -o prometheus-p1-exporter -ldflags "-X main.GitCommit=${GIT_COMMIT}"

FROM alpine:latest
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=builder /build/prometheus-p1-exporter /app
EXPOSE 8888
CMD ["/app/prometheus-p1-exporter", "-listen", "0.0.0.0:8888"]
