FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN go build -o ./program ./cmd/server/main.go

FROM alpine:latest
WORKDIR /root
COPY --from=builder /app/program ./
COPY --from=builder /app/migrations ./migrations
CMD ["./program"]