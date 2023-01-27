FROM golang:1.19.4 AS builder
ENV GO111MODULE=on

WORKDIR /go/src/github.com/thavlik/transcriber
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN cd cmd && go build -o transcriber

FROM debian:bullseye-slim
COPY --from=builder /go/src/github.com/thavlik/transcriber/cmd/transcriber /usr/local/bin
CMD ["/usr/local/bin/transcriber"]
