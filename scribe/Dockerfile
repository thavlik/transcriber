ARG FDK_AAC_VERSION=2.0.0
ARG KINDEST_REPOSITORY=""
ARG KINDEST_TAG="latest"
ARG RUNTIME_BASE_IMAGE="debian:bullseye-slim"
FROM ${KINDEST_REPOSITORY}thavlik/transcriber-base:${KINDEST_TAG} AS builder
COPY base/pkg base/pkg
COPY comprehend/pkg/comprehend comprehend/pkg/comprehend
COPY scribe scribe
RUN cd scribe/cmd && go build -o scribe
FROM ${RUNTIME_BASE_IMAGE}
RUN apt-get update \
    && apt-get install -y \
        ca-certificates \
    && apt-get clean
ARG FDK_AAC_VERSION
COPY --from=builder /usr/local/fdk-aac-${FDK_AAC_VERSION}/lib/libfdk-aac.so /usr/local/lib
COPY --from=builder /usr/local/fdk-aac-${FDK_AAC_VERSION}/lib/libfdk-aac.so.2 /usr/local/lib
COPY --from=builder /usr/local/fdk-aac-${FDK_AAC_VERSION}/lib/libfdk-aac.so.2.0.0 /usr/local/lib
COPY --from=builder /go/src/github.com/thavlik/transcriber/scribe/cmd/scribe /usr/local/bin
ENV LD_LIBRARY_PATH=/usr/local/lib
CMD ["/usr/local/bin/scribe"]
