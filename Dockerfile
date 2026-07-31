# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o expensesync main.go

# Production runtime stage (pinned alpine version for reproducible builds)
FROM alpine:3.20

RUN apk --no-cache add ca-certificates

WORKDIR /app/

# Copy binary from builder stage
COPY --from=builder /app/expensesync .

EXPOSE 8080

CMD ["./expensesync"]
