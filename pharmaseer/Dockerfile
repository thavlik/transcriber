ARG KINDEST_REPOSITORY=""
ARG KINDEST_TAG="latest"
ARG RUNTIME_BASE_IMAGE="ghcr.io/puppeteer/puppeteer:latest"
FROM ${KINDEST_REPOSITORY}thavlik/transcriber-base:${KINDEST_TAG} AS builder
COPY base/pkg base/pkg
COPY pharmaseer pharmaseer
WORKDIR /go/src/github.com/thavlik/transcriber/pharmaseer/cmd
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o pharmaseer
FROM ${RUNTIME_BASE_IMAGE}
USER root
COPY --from=builder /go/src/github.com/thavlik/transcriber/pharmaseer/cmd/pharmaseer /usr/local/bin
COPY pharmaseer/scripts /scripts
CMD ["pharmaseer"]

