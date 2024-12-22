# Use the official Node.js image based on Debian Bullseye
FROM node:23.5.0-bullseye-slim

# Update package list and install git
RUN apt-get update && apt-get install -y git

WORKDIR /home/app

# Copy the built Go binary from your local machine
COPY ./main .
COPY .env .

# Ensure the binary has execute permissions
RUN chmod +x ./main

# Command to run the executable
ENTRYPOINT ["./main"]