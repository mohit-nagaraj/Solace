# Use the official Golang image to create a build artifact.
FROM golang:1.21 as builder

# Set the Current Working Directory inside the container
WORKDIR /home/app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .

# MULTI STAGE BUILDS

# Use the official Node.js image to install Node.js dependencies
FROM node:14

RUN apt-get update && apt-get install git -y

WORKDIR /home/app

# Copy the built Go binary from the builder stage
COPY --from=builder /home/app/main .

# Ensure the binary has execute permissions
RUN chmod +x ./main

# Command to run the executable
CMD ["./main"]