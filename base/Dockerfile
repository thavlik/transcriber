ARG FDK_AAC_VERSION=2.0.0
FROM debian:bullseye-slim AS fdkaac
RUN apt-get update \
    && apt-get install -y \
        wget \
        tar \
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
FROM golang:1.19.4
ARG FDK_AAC_VERSION
COPY --from=fdkaac /usr/local/fdk-aac-${FDK_AAC_VERSION} /usr/local/fdk-aac-${FDK_AAC_VERSION}
ENV CGO_CPPFLAGS="-I/usr/local/fdk-aac-${FDK_AAC_VERSION}/include/fdk-aac"
ENV CGO_LDFLAGS="-L/usr/local/fdk-aac-${FDK_AAC_VERSION}/lib"
ENV GO111MODULE=on
WORKDIR /go/src/github.com/thavlik/transcriber
COPY go.mod .
COPY go.sum .
RUN go mod download
