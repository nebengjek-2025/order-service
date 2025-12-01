FROM golang:1.25-alpine AS builder

RUN apk add --no-cache \
    build-base \
    ca-certificates \
    git \
    librdkafka-dev \
    pkgconfig \
    tzdata 

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY src ./src

ENV CGO_ENABLED=1


RUN go build -tags dynamic -ldflags="-s -w" -o order-service ./src/cmd/app/main.go

FROM alpine:3.20

RUN apk add --no-cache \
    ca-certificates \
    librdkafka \
    tzdata && \
    adduser -D -g '' appuser


WORKDIR /app

COPY --from=builder /app/order-service ./order-service

ENV TZ=Asia/Jakarta

EXPOSE 9000

USER appuser

ENTRYPOINT ["./order-service"]
