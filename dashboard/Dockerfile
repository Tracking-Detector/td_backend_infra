FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY ./dashboard /app
RUN go mod download

RUN go build -o main ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/static ./static

COPY --from=builder /app/main .

CMD ["./main"]