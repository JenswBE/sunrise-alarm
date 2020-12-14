ARG SERVICE_NAME=srv-alarm

# Setup builder
FROM --platform=${BUILDPLATFORM} rust AS builder-base

FROM builder-base AS builder-amd64
ENV TARGET=x86_64-unknown-linux-musl
RUN apt-get update && \
    apt-get install -qq \
    build-essential \
    musl-tools

FROM builder-base AS builder-arm
ENV TARGET=armv7-unknown-linux-musleabihf
ENV CC_armv7_unknown_linux_musleabihf=arm-linux-gnueabihf-gcc
RUN apt-get update && \
    apt-get install -qq \
    crossbuild-essential-armhf

# Build service
FROM builder-${TARGETARCH} AS builder
ARG SERVICE_NAME
WORKDIR /usr/src/sunrise-alarm
RUN rustup target add ${TARGET}

COPY . .
WORKDIR /usr/src/sunrise-alarm/${SERVICE_NAME}
RUN cargo test
RUN cargo build --target ${TARGET} --release 
RUN cp /usr/src/sunrise-alarm/target/${TARGET}/release/${SERVICE_NAME} /service

# Build final image
FROM scratch
ENV WARP_PORT 8080
ENV HOST_SRV_ALARM http://srv-alarm:8080
ENV HOST_SRV_CONFIG http://srv-config:8080
ENV HOST_SRV_PHYSICAL http://srv-physical:8080
ENV HOST_SRV_AUDIO http://srv-audio:8080
EXPOSE 8080

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /service service
CMD ["./service"]