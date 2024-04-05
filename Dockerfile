# Stage 1: Build the Go binary
FROM golang:1.22 AS builder

WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go binary
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
RUN go build -o microservice .

# Stage 2: Create a minimal image to run the binary
FROM alpine:3.19

# Create a non-root user
RUN adduser -D -u 1000 appuser

RUN apk --no-cache add ca-certificates

WORKDIR /home/appuser

# Copy the binary from the builder stage
COPY --from=builder /app/microservice .

# Set ownership and permissions for the non-root user
RUN chown appuser:appuser microservice && \
    chmod 755 microservice

USER appuser

EXPOSE 3000

# Command to run the binary
CMD ["./microservice"]
