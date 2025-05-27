FROM golang:1.23 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/server ./cmd/server

FROM alpine:3.17
WORKDIR /app
COPY --from=builder /app/server .
COPY .env .
ENTRYPOINT ["/app/server"]

   