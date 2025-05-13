FROM golang:1.23-alpine3.21 AS builder
WORKDIR /app
COPY . .
WORKDIR /app/cmd/ordersystem
RUN go build -o /app/main

FROM alpine:3.21
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]