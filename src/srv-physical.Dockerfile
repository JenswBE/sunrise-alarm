ARG SERVICE_NAME=srv-physical

FROM --platform=${BUILDPLATFORM} golang:1.16 AS builder-base

FROM builder-base AS builder-amd64
ENV GOOS=linux
ENV GOARCH=amd64

FROM builder-base AS builder-arm
ENV GOOS=linux
ENV GOARCH=arm
ENV GOARM=7

FROM builder-${TARGETARCH} AS builder
ARG SERVICE_NAME
WORKDIR /src/
COPY ${SERVICE_NAME} .
WORKDIR /src/api
RUN CGO_ENABLED=0 go build -ldflags='-extldflags=-static' -o /bin/srv-physical

FROM alpine
COPY --from=builder /bin/srv-physical /srv-physical/bin/srv-physical
COPY --from=builder /src/docs/index.html /srv-physical/docs/index.html
COPY --from=builder /src/docs/oauth2-redirect.html /srv-physical/docs/oauth2-redirect.html
COPY --from=builder /src/docs/openapi.yml /srv-physical/docs/openapi.yml
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080
WORKDIR /srv-physical/bin
ENTRYPOINT ["./srv-physical"]