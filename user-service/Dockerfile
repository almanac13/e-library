FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o user-service ./cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/user-service .
COPY .env .

EXPOSE 8087
EXPOSE 50051

CMD ["./user-service"]