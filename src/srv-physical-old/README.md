# srv-physical

## Run locally
```bash
# Prepare
python3 -m venv env
source env/bin/activate
pip install -U pip wheel
pip install -r requirements.txt

# Run
Mock=True gunicorn physical.main:app -c gunicorn.conf.py -b localhost:8002 --reload
```
