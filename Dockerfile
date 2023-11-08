# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
FROM golang:1.21 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN GOOS=linux go build -o matching ./cmd/tinder-server/main.go

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./matching"]
