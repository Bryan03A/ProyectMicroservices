# Use an official Go image as the base
FROM golang:1.19 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the service files into the container
COPY . .

# Download dependencies
RUN go mod tidy

# Compile the application binary
RUN go build -o orders-service

# Second stage: a lighter image to run the application
FROM debian:bullseye-slim

# Set the working directory
WORKDIR /app

# Copy the binary from the first stage
COPY --from=builder /app/orders-service .

# Expose the port where the service will run
EXPOSE 5008

# Command to run the application
CMD ["./orders-service"]