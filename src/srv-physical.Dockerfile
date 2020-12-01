ARG SERVICE_NAME=srv-physical

FROM python:3.8-slim
ARG SERVICE_NAME

# Install dependencies
COPY ${SERVICE_NAME}/requirements.txt .
RUN apt-get update && apt-get -qq install build-essential
RUN pip install --no-cache-dir -U pip wheel && \
    pip install --no-cache-dir gunicorn uvicorn uvloop httptools
RUN pip install --no-cache-dir -r requirements.txt

# Copy service
COPY ${SERVICE_NAME} .

# Limiting workers is required for MQTT to work correctly
EXPOSE 8080
CMD [ "gunicorn", "physical.main:app", "-w", "1", "-k", "uvicorn.workers.UvicornWorker", "-b", "0.0.0.0:8080" ]