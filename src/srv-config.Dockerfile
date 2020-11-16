ARG SERVICE_NAME=srv-config
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
    cargo

# Build service
WORKDIR /home/rust/${SERVICE_NAME}
RUN cargo test
RUN cargo build --release

# Start building the final image
FROM ${BASE_IMAGE}
ARG SERVICE_NAME
ENV WARP_PORT 80
ENV DATA_DIR_PATH /data
EXPOSE 80
VOLUME [ "/data" ]

RUN apk add --no-cache libgcc
COPY --from=builder /home/rust/target/release/${SERVICE_NAME} service
CMD ["./service"]