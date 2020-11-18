# Based on https://vuejs.org/v2/cookbook/dockerize-vuejs-app.html

ARG SERVICE_NAME=gui-watchface

# Build application
FROM node:lts-alpine as builder
ARG SERVICE_NAME
WORKDIR /app
COPY ${SERVICE_NAME}/package.json ./
RUN yarn install
COPY ${SERVICE_NAME} .
RUN yarn build

# Build final image
FROM nginx:stable-alpine
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]