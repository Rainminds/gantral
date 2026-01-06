
FROM python:3.11-slim

WORKDIR /app

# Install dependencies
RUN pip install requests

# Copy runner code
COPY runner.py .

# The agent code will be mounted at runtime to /agent
CMD ["python", "-u", "runner.py"]
