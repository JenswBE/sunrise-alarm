ARG SERVICE_NAME=srv-physical

# Setup base
FROM --platform=${TARGETPLATFORM} python:3.8-slim AS base

# Install OS dependencies
FROM base AS base-amd64
ENV MOCK True
RUN apt-get update && apt-get -qq install build-essential

FROM base AS base-arm
ENV MOCK False
RUN apt-get update && apt-get -qq install build-essential python3-rpi.gpio libgpiod2
RUN pip install --no-cache-dir -U RPi.GPIO

# Install python dependencies
FROM base-${TARGETARCH}
ARG SERVICE_NAME
COPY ${SERVICE_NAME}/requirements.txt .
RUN pip install --no-cache-dir -U pip wheel && \
    pip install --no-cache-dir gunicorn uvicorn uvloop httptools
RUN CFLAGS="-fcommon" pip install --no-cache-dir -r requirements.txt # See https://forum.manjaro.org/t/pip-install-rpi-gpio-fail/25788/3

# Copy service
COPY ${SERVICE_NAME} .

# Limiting workers is required for MQTT to work correctly
EXPOSE 8080
CMD [ "gunicorn", "physical.main:app", "-c", "gunicorn.conf.py" ]