# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /telegram-notifier

# Final stage
FROM alpine:latest
WORKDIR /
COPY --from=builder /telegram-notifier /telegram-notifier
EXPOSE 8080
ENTRYPOINT ["/telegram-notifier"]
