ARG SERVICE_NAME=srv-audio
ARG BASE_IMAGE=alpine:3.12

# Copy sources
FROM ${BASE_IMAGE} as builder
ARG SERVICE_NAME
WORKDIR /home/rust/
COPY . .

# Install requirements
RUN apk add --no-cache \
    gcc \
    musl-dev \
    rust \
    cargo \
    alsa-lib-dev

# Build service
WORKDIR /home/rust/${SERVICE_NAME}
RUN cargo test
RUN cargo build --release

# Start building the final image
FROM ${BASE_IMAGE}
ARG SERVICE_NAME
ENV WARP_PORT 80
ENV MUSIC_DIR_PATH /music
EXPOSE 80
VOLUME [ "/music" ]

# Setup audio
# Based on https://wiki.alpinelinux.org/wiki/Sound_Setup
RUN apk add --no-cache \
    alsa-utils \
    alsa-lib \
    alsaconf \
    libgcc \
    pulseaudio \
    pulseaudio-alsa \
    alsa-plugins-pulse

COPY ${SERVICE_NAME}/default.mp3 .
COPY ${SERVICE_NAME}/pulseaudio_client.conf /etc/pulse/client.conf
COPY --from=builder /home/rust/target/release/${SERVICE_NAME} service
CMD ["./service"]