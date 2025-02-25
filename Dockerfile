FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /exchanger .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /exchanger .
EXPOSE 8080

CMD ["./exchanger"]