FROM golang:alpine
WORKDIR $GOPATH/src/app/
ADD . .
RUN apk add --no-cache git make bash
RUN make
RUN make install
ENTRYPOINT ["prometheus-p1-exporter", "-listen", "0.0.0.0:8888"]