
FROM python:3.11-slim

WORKDIR /app

RUN pip install requests

COPY runner.py .

CMD ["python", "-u", "runner.py"]
