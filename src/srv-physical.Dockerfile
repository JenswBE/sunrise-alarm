ARG SERVICE_NAME=srv-physical

# Setup base
FROM --platform=${TARGETPLATFORM} python:3.8-slim AS base

# Install OS dependencies
FROM base AS base-amd64
RUN apt-get update && apt-get -qq install build-essential

FROM base AS base-arm
RUN apt-get update && apt-get -qq install build-essential python3-rpi.gpio

# Install python dependencies
FROM builder-${TARGETARCH}
ARG SERVICE_NAME
COPY ${SERVICE_NAME}/requirements.txt .
RUN pip install --no-cache-dir -U pip wheel && \
    pip install --no-cache-dir gunicorn uvicorn uvloop httptools
RUN pip install --no-cache-dir -r requirements.txt

# Copy service
COPY ${SERVICE_NAME} .

# Limiting workers is required for MQTT to work correctly
EXPOSE 8080
CMD [ "gunicorn", "physical.main:app", "-w", "1", "-k", "uvicorn.workers.UvicornWorker", "-b", "0.0.0.0:8080" ]