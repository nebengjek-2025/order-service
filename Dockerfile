# =======================
# 1. Builder Stage
# =======================
FROM golang:1.25-alpine AS builder

# Install tools build + librdkafka dev (buat confluent-kafka-go)
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    build-base \
    librdkafka-dev \
    pkgconfig

WORKDIR /app

# Copy go.mod & go.sum dulu supaya cache oke
COPY go.mod go.sum ./
RUN go mod download

# Copy semua source
COPY . .

# Enable CGO (wajib untuk confluent-kafka-go)
ENV CGO_ENABLED=1

# ❗ PENTING: pakai build tag "dynamic" supaya pakai librdkafka dari sistem,
# bukan librdkafka_glibc_linux.a yang dibundel di module (yang bikin error tadi)
RUN go build -tags dynamic -ldflags="-s -w" -o order-service ./src/cmd/app/main.go

# =======================
# 2. Runtime Stage
# =======================
FROM alpine:3.20

RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    librdkafka && \
    adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder /app/order-service ./order-service

ENV TZ=Asia/Jakarta

EXPOSE 8080

USER appuser

ENTRYPOINT ["./order-service"]
