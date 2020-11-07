FROM tiangolo/uvicorn-gunicorn-fastapi:python3.8

# Setup ENV
# Limiting workers is required for MQTT to work correctly
ENV MODULE_NAME physical.main
ENV MAX_WORKERS 1
ENV WEB_CONCURRENCY 1

# Copy srv-physical into container
COPY . /app/

# Install dependencies
RUN pip install -U pip wheel
RUN pip install -r /app/requirements.txt