ARG KINDEST_REPOSITORY=""
ARG KINDEST_TAG="latest"
ARG RUNTIME_BASE_IMAGE="debian:bullseye-slim"
FROM ${KINDEST_REPOSITORY}thavlik/transcriber-base:${KINDEST_TAG} AS builder
COPY base/pkg base/pkg
COPY define define
WORKDIR /go/src/github.com/thavlik/transcriber/define/cmd
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o define
FROM ${RUNTIME_BASE_IMAGE}
RUN apt-get update \
    && apt-get install -y \
        ca-certificates \
    && apt-get clean
COPY --from=builder /go/src/github.com/thavlik/transcriber/define/cmd/define /usr/local/bin
CMD ["define"]
