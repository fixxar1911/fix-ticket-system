FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o ticket-service .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/ticket-service .

EXPOSE 8080

CMD ["./ticket-service"] 