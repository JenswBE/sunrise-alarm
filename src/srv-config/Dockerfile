# Based on https://blog.sedrik.se/posts/my-docker-setup-for-rust/
FROM ekidd/rust-musl-builder:stable as builder

WORKDIR /home/rust/

# Avoid having to install/build all dependencies by copying
# the Cargo files and making a dummy src/main.rs
COPY Cargo.toml .
COPY Cargo.lock .
RUN echo "fn main() {}" > src/main.rs
RUN cargo test
RUN cargo build --release

# We need to touch our real main.rs file or else docker will use
# the cached one.
COPY . .
RUN sudo touch src/main.rs

RUN cargo test
RUN cargo build --release

# Start building the final image
FROM scratch
WORKDIR /home/rust/
VOLUME [ "/data" ]
COPY --from=builder /home/rust/target/x86_64-unknown-linux-musl/release/srv-config .

ENV WARP_PORT 80
ENV DATA_DIR_PATH /data
EXPOSE 80
ENTRYPOINT ["./srv-config"]