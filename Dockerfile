FROM golang:1.25-alpine as builder

# Install needed tools
RUN apk add --no-cache bash git openssh curl vim busybox-extras procps

ADD ./ /app
WORKDIR /app/cmd

RUN go build -o order-service .

RUN cd /app

# Run stage
FROM --platform=linux/amd64 debian:bullseye-slim

# Update package lists and upgrade packages to fix vulnerabilities
RUN apt-get update && apt-get upgrade -y

# Set working directory
WORKDIR /app/cmd
COPY --from=builder /app/cmd/utility-monitoring .

RUN apt-get install -y tzdata && apt-get clean

CMD ["/app/cmd/order-service"]