# Stage 1: Build
FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/main.go

# Stage 2: Run
FROM ubuntu:latest

WORKDIR /app

# Устанавливаем корневые сертификаты и curl
RUN apt-get update && \
    apt-get install -y \
        ca-certificates \
        curl && \
    update-ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/main .

EXPOSE 5000

CMD ["./main"]
