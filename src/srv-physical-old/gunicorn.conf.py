bind = "0.0.0.0:8080"
logconfig = "logging.conf"
accesslog = "-"
workers = 1
worker_class = "uvicorn.workers.UvicornWorker"
