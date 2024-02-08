FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY ./microservices/go.mod /app/go.mod
COPY ./microservices/go.sum /app/go.sum
RUN go mod download

COPY ./microservices/functions/model/main.go /app/main.go
COPY ./microservices/shared /app/shared

# Build the function executable with optimizations enabled
RUN go build -o main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

CMD ["./main"]