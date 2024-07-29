# ---- Base Stage ----
FROM golang:1.22-alpine AS base

# Set necessary environment variables for Go modules
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set the working directory inside the container
WORKDIR /restaurant-web

# Install git
RUN apk update && apk add --no-cache ca-certificates git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# ---- Build Stage ----
FROM base AS build

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o ./bin/server ./cmd

# ---- Production Stage ----
FROM alpine:latest AS production

# Set the working directory inside the container
WORKDIR /restaurant-web

# Create a non-root user
RUN adduser -D -u 1000 appuser

# Copy the pre-built binary from the build stage
COPY --from=build /restaurant-web/bin/server .

# Copy the .env file into the container
COPY .env .env

# Set ownership of the application directory to the non-root user
RUN chown -R appuser: /restaurant-web

# Switch to the non-root user
USER appuser

# Expose port 3020 to the outside world
EXPOSE 3021

# Command to run the executable
CMD ["./server"]
