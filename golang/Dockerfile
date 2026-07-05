# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o tyk-sre-assignment .

# Final minimal stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/tyk-sre-assignment .
EXPOSE 8080
CMD ["./tyk-sre-assignment"]