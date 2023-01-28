FROM debian:bullseye AS builder
ARG FLUTTER_VERSION=3.3.9
RUN apt-get update \
    && apt-get install -y \
        curl \
        git \
        wget \
        zip \
        unzip \
        libgconf-2-4 \
        gdb \
        libstdc++6 \
        libglu1-mesa \
        fonts-droid-fallback \
        lib32stdc++6 \
        python3 \
    && apt-get clean
ENV PATH="/flutter/bin:/flutter/bin/cache/dart-sdk/bin:${PATH}"
#RUN git clone -q \
#    -c advice.detachedHead=false \
#    --depth 1 \
#    --branch ${FLUTTER_VERSION} \
#    https://github.com/flutter/flutter.git \
#    /flutter \
#    && flutter config --enable-web \
#    && flutter doctor -v
RUN git clone -q \
    https://github.com/flutter/flutter.git \
    /flutter \
    && flutter channel master \
    && flutter config --enable-web \
    && flutter upgrade \
    && flutter doctor -v
WORKDIR /app
COPY pubspec.lock .
COPY pubspec.yaml .
RUN flutter pub get

COPY . .
RUN flutter build web

FROM nginx:1.23
COPY --from=builder /app/build/web /etc/nginx/html
