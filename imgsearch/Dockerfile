ARG KINDEST_REPOSITORY=""
ARG KINDEST_TAG="latest"
ARG RUNTIME_BASE_IMAGE="debian:bullseye-slim"
FROM ${KINDEST_REPOSITORY}thavlik/transcriber-base:${KINDEST_TAG} AS builder
COPY base/pkg base/pkg
COPY imgsearch imgsearch
WORKDIR /go/src/github.com/thavlik/transcriber/imgsearch/cmd
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o imgsearch
FROM ${RUNTIME_BASE_IMAGE}
RUN apt-get update \
    && apt-get install -y ca-certificates \
    && rm -rf /var/lib/apt/lists/*
COPY --from=builder /go/src/github.com/thavlik/transcriber/imgsearch/cmd/imgsearch /usr/local/bin
CMD ["imgsearch"]
