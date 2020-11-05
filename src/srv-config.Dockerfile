FROM ekidd/rust-musl-builder:stable as builder

ARG SERVICE_NAME=srv-config

# Copy sources
WORKDIR /home/rust/
COPY . .

# Build service
WORKDIR /home/rust/${SERVICE_NAME}
RUN cargo test
RUN cargo build --release

# Start building the final image
FROM scratch

ARG SERVICE_NAME=srv-config
ENV WARP_PORT 80
ENV DATA_DIR_PATH /data
EXPOSE 80
VOLUME [ "/data" ]

COPY --from=builder /home/rust/target/x86_64-unknown-linux-musl/release/${SERVICE_NAME} service
ENTRYPOINT ["./service"]