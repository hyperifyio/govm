# Build stage
FROM golang:1.22 AS builder
WORKDIR /app

# Cache and install dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY Makefile .
COPY doc.go .
COPY cmd ./cmd
COPY internal ./internal

# Build the application
RUN make build

# Final stage
FROM alpine:latest
WORKDIR /app/

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/govm .

EXPOSE 8080

# Command to run
CMD ["./govm"]
