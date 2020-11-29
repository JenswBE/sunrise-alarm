ARG SERVICE_NAME=srv-audio
ARG BASE_IMAGE=debian:buster-slim

# Setup builder
FROM --platform=${BUILDPLATFORM} rust AS builder-base

FROM builder-base AS builder-amd64
ENV TARGET=x86_64-unknown-linux-gnu
RUN apt-get update && \
    apt-get install -qq \
    build-essential \
    libasound2-dev

FROM builder-base AS builder-arm
ENV TARGET=armv7-unknown-linux-gnueabihf
ENV PKG_CONFIG_LIBDIR_armv7_unknown_linux_gnueabihf=/usr/lib/arm-linux-gnueabihf/pkgconfig
RUN dpkg --add-architecture armhf && \
    apt-get update && \
    apt-get install -qq \
    crossbuild-essential-armhf \
    libasound2-dev \
    libasound2-dev:armhf

# Build service
FROM builder-${TARGETARCH} AS builder
ARG SERVICE_NAME
ENV PKG_CONFIG_ALLOW_CROSS=true
WORKDIR /usr/src/sunrise-alarm
COPY . .

RUN rustup target add ${TARGET}
WORKDIR /usr/src/sunrise-alarm/${SERVICE_NAME}
RUN cargo test
RUN cargo build --target ${TARGET} --release
RUN cp /usr/src/sunrise-alarm/target/${TARGET}/release/${SERVICE_NAME} /service

# Build final image
FROM ${BASE_IMAGE}
ARG SERVICE_NAME
ENV WARP_PORT 8080
ENV HOST_SRV_ALARM srv-alarm:8080
ENV HOST_SRV_CONFIG srv-config:8080
ENV HOST_SRV_PHYSICAL srv-physical:8080
ENV HOST_SRV_AUDIO srv-audio:8080
ENV MUSIC_DIR_PATH /music
EXPOSE 8080
VOLUME [ "/music" ]

RUN apt-get update && \
    apt-get install -y \
    libasound2 \
    pulseaudio \
    && \
    rm -rf /var/lib/apt/lists/*

COPY ${SERVICE_NAME}/default.mp3 .
COPY ${SERVICE_NAME}/pulseaudio_client.conf /etc/pulse/client.conf
COPY --from=builder /service service
CMD ["./service"]