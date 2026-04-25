FROM golang:1.25.3-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o main ./cmd/api/main.go

FROM alpine:latest

COPY --from=builder /app /

CMD ["./main"]