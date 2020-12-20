# srv-physical

## Run locally
```bash
# Prepare
python3 -m venv env
source env/bin/activate
pip install -U pip wheel
pip install -r requirements.txt

# Run
export Mock=True
gunicorn physical.main:app -w 1 -k uvicorn.workers.UvicornWorker -b localhost:8002 --reload --log-config=logging.conf
```
