# Based on https://github.com/nuxt/nuxt.js/blob/dev/examples/docker-build/Dockerfile

ARG SERVICE_NAME=gui-watchface

# Setup builder
FROM node:lts-alpine as builder
ARG SERVICE_NAME
WORKDIR /src
COPY ${SERVICE_NAME} .

# Build application
RUN yarn install --immutable
RUN yarn build

# Only install Production dependencies
RUN rm -rf node_modules
RUN NODE_ENV=production yarn workspaces focus --all --production
RUN yarn cache clean --all

# Build final image
FROM node:lts-alpine
WORKDIR /src
COPY --from=builder /src  .
EXPOSE 8080
CMD [ "yarn", "start", "--hostname", "0.0.0.0", "--port", "8080" ]
