FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o notification-service ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/notification-service .
COPY .env .

CMD ["./notification-service"]