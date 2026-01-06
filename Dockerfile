
FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build Server
RUN go build -o /app/bin/server ./cmd/server
# Build Worker
RUN go build -o /app/bin/worker ./cmd/worker

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/bin/server /app/server
COPY --from=builder /app/bin/worker /app/worker

# Install basic tools
RUN apk add --no-cache bash curl

# Create start script
RUN echo '#!/bin/bash' > /app/start.sh && \
    echo 'set -e' >> /app/start.sh && \
    echo '/app/worker &' >> /app/start.sh && \
    echo 'WORKER_PID=$!' >> /app/start.sh && \
    echo '/app/server &' >> /app/start.sh && \
    echo 'SERVER_PID=$!' >> /app/start.sh && \
    echo 'wait -n' >> /app/start.sh && \
    echo 'exit $?' >> /app/start.sh && \
    chmod +x /app/start.sh

CMD ["/app/start.sh"]
