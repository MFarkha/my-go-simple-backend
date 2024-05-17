# Stage 1: Build the Go binary
FROM golang:1.22 AS builder

WORKDIR /app

# # Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY *.go .

# Build the Go binary - with statically linked library to run on Alpine Linux
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o microservice .
RUN CGO_ENABLED=0 GOOS=linux go build -o microservice .

# Stage 2: Create a minimal image to run the binary
FROM alpine:3.19

# Create a non-root user
RUN adduser -D -u 1000 appuser

RUN apk --no-cache add ca-certificates

WORKDIR /home/appuser

# Copy the binary from the builder stage
COPY --from=builder /app/microservice ./

# Copy the dotenv file
COPY .env ./

# Set ownership and permissions for the non-root user
RUN chown appuser:appuser microservice && \
    chmod 755 microservice

USER appuser

# ARG PORT
# ARG MAX_RANDOM_NUMBER
# ARG METRIC_DECIMAL_PLACES

# ENV PORT=${PORT} \
#     APP_MAX_RANDOM_NUMBER=${MAX_RANDOM_NUMBER} \
#     APP_METRIC_DECIMAL_PLACES=${METRIC_DECIMAL_PLACES}
# EXPOSE ${PORT}

ENTRYPOINT ["./microservice"]
