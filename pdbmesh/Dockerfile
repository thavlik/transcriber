ARG KINDEST_REPOSITORY=""
ARG KINDEST_TAG="latest"
ARG RUNTIME_BASE_IMAGE="debian:bullseye-slim"
FROM ${KINDEST_REPOSITORY}thavlik/transcriber-base:${KINDEST_TAG} AS builder
COPY base/pkg base/pkg
COPY pdbmesh pdbmesh
WORKDIR /go/src/github.com/thavlik/transcriber/pdbmesh/cmd
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o pdbmesh
FROM ${RUNTIME_BASE_IMAGE}
RUN apt-get update \
    && apt-get install -y ca-certificates \
    && rm -rf /var/lib/apt/lists/*
COPY --from=builder /go/src/github.com/thavlik/transcriber/pdbmesh/cmd/pdbmesh /usr/local/bin
CMD ["pdbmesh"]
