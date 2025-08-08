FROM golang:1.21 AS builder
WORKDIR /app
COPY ../Homework_5/go.mod ./
COPY ../Homework_5/internal/ ./internal/
COPY ../Homework_5/cmd/ ./cmd/
RUN go build -o app ./cmd/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/app .
CMD ["./app"]
