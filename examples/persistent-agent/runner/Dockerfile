
FROM python:3.11-slim

WORKDIR /app

# Install dependencies
RUN pip install requests pyjwt

# Copy runner code
COPY . .

# The agent code will be mounted at runtime to /agent
CMD ["python", "-u", "runner.py"]
