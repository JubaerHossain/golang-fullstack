# Start from a minimal Go image
FROM golang:1.22-alpine AS build

# Set necessary environment variables required for Go
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set the working directory inside the container
WORKDIR /app

# Install git and ca-certificates
RUN apk update && apk add --no-cache ca-certificates git

# Use GitHub as the primary Go module proxy
ENV GOPROXY=github.com,direct

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire source code into the container
COPY . .

# Build the Go app
RUN go build -o app cmd/main.go

# Start a new stage from scratch
FROM alpine:latest  

# Install necessary packages (e.g., for running the app)
RUN apk --no-cache add ca-certificates

# Copy the pre-built binary from the previous stage
COPY --from=build /app/app /usr/local/bin/app
# Copy the .env file into the container
COPY .env .env

# Expose port 3031 (change to your application's port if needed)
EXPOSE 3031

# Command to run the executable
CMD ["app"]
