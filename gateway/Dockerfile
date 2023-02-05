ARG KINDEST_REPOSITORY=""
ARG KINDEST_TAG="latest"
ARG RUNTIME_BASE_IMAGE="alpine:3.17.0"
FROM ${KINDEST_REPOSITORY}thavlik/transcriber-base:${KINDEST_TAG} AS builder
COPY base/cmd/iam base/cmd/iam
COPY base/pkg base/pkg
COPY pharmaseer/pkg/api pharmaseer/pkg/api
COPY gateway gateway
WORKDIR /go/src/github.com/thavlik/transcriber/gateway/cmd
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o gateway
FROM ${RUNTIME_BASE_IMAGE}
COPY --from=builder /go/src/github.com/thavlik/transcriber/gateway/cmd/gateway /usr/local/bin
CMD ["gateway"]
