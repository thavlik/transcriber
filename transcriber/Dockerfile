ARG FDK_AAC_VERSION=2.0.0
ARG KINDEST_REPOSITORY=""
ARG KINDEST_TAG="latest"
ARG RUNTIME_BASE_IMAGE="debian:bullseye-slim"
FROM ${KINDEST_REPOSITORY}thavlik/transcriber-base:${KINDEST_TAG} AS builder

RUN apt-get update \
    && apt-get install -y \
        build-essential \
        libtool \
        autoconf \
        automake \
        autotools-dev \
    && apt-get clean

ARG FDK_AAC_VERSION
RUN mkdir -p /tmp/fdk-aac \
    && cd /tmp/fdk-aac \
    && wget https://github.com/mstorsjo/fdk-aac/archive/refs/tags/v${FDK_AAC_VERSION}.tar.gz \
    && tar -zxvf v${FDK_AAC_VERSION}.tar.gz \
    && cd fdk-aac-${FDK_AAC_VERSION}/ \
    && ./autogen.sh \
    && ./configure --prefix=/usr/local/fdk-aac-${FDK_AAC_VERSION} \
    && make \
    && make install \
    && cd / \
    && rm -rf /tmp/fdk-aac

ENV GO111MODULE=on
WORKDIR /go/src/github.com/thavlik/transcriber
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
ENV CGO_CPPFLAGS="-I/usr/local/fdk-aac-${FDK_AAC_VERSION}/include/fdk-aac"
ENV CGO_LDFLAGS="-L/usr/local/fdk-aac-${FDK_AAC_VERSION}/lib"
RUN cd cmd && go build -o transcriber

FROM ${RUNTIME_BASE_IMAGE}
RUN apt-get update \
    && apt-get install -y \
        ca-certificates \
    && apt-get clean
ARG FDK_AAC_VERSION
COPY --from=builder /usr/local/fdk-aac-${FDK_AAC_VERSION}/lib/libfdk-aac.so /usr/local/lib
COPY --from=builder /usr/local/fdk-aac-${FDK_AAC_VERSION}/lib/libfdk-aac.so.2 /usr/local/lib
COPY --from=builder /usr/local/fdk-aac-${FDK_AAC_VERSION}/lib/libfdk-aac.so.2.0.0 /usr/local/lib
COPY --from=builder /go/src/github.com/thavlik/transcriber/cmd/transcriber /usr/local/bin
ENV LD_LIBRARY_PATH=/usr/local/lib
CMD ["/usr/local/bin/transcriber"]
