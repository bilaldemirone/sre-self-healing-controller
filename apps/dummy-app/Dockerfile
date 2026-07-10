FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY main.go .

RUN go build -o dummy-app main.go

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/dummy-app .

EXPOSE 8080

CMD ["./dummy-app"]