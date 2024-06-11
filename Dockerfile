# Use the official Golang image to build and run the Go application
FROM golang:1.22 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Copy the .env file
COPY .env .env

# Build the Go app
RUN go build -o main .

# Use the same base image for the final stage
FROM golang:1.22

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Pre-built binary file and .env from the previous stage
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
