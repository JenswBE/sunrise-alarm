ARG SERVICE_NAME=srv-config

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
ENV CARGO_TARGET_ARMV7_UNKNOWN_LINUX_MUSLEABIHF_LINKER=arm-linux-gnueabihf-gcc
ENV CC_armv7_unknown_linux_musleabihf=arm-linux-gnueabihf-gcc
RUN apt-get update && \
    apt-get install -qq \
    crossbuild-essential-armhf

# Build service
FROM builder-${TARGETARCH} AS builder
ARG SERVICE_NAME
WORKDIR /usr/src/sunrise-alarm
COPY . .

RUN rustup target add ${TARGET}
WORKDIR /usr/src/sunrise-alarm/${SERVICE_NAME}
RUN cargo test
RUN cargo build --target ${TARGET} --release 
RUN cp /usr/src/sunrise-alarm/target/${TARGET}/release/${SERVICE_NAME} /service

# Build final image
FROM scratch
ARG SERVICE_NAME
ENV WARP_PORT 80
ENV DATA_DIR_PATH /data
EXPOSE 80
VOLUME [ "/data" ]

COPY --from=builder /service service
CMD ["./service"]