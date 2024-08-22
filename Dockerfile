FROM golang:1.21-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

# Download all dependencies.
RUN go mod download

COPY . .

# Build the Go app
RUN go build -o my-pokemon-cache .

# Start a new stage from scratch
FROM alpine:latest

# Copy the Pre-built binary file from the previous stage
COPY --from=build /app/my-pokemon-cache /usr/local/bin/my-pokemon-cache

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["my-pokemon-cache"]
