ARG KINDEST_REPOSITORY=""
ARG KINDEST_TAG="latest"
ARG RUNTIME_BASE_IMAGE="debian:bullseye-slim"
FROM ${KINDEST_REPOSITORY}thavlik/transcriber-base:${KINDEST_TAG} AS builder
COPY base/pkg base/pkg
COPY broadcaster broadcaster
WORKDIR /go/src/github.com/thavlik/transcriber/broadcaster/cmd
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o broadcaster
FROM ${RUNTIME_BASE_IMAGE}
RUN apt-get update \
    && apt-get install -y \
        ca-certificates \
    && apt-get clean
COPY --from=builder /go/src/github.com/thavlik/transcriber/broadcaster/cmd/broadcaster /usr/local/bin
CMD ["broadcaster"]
