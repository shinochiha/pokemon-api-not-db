# Use an official Golang runtime as the base image
FROM golang:alpine

# Set the working directory in the container
WORKDIR /app

# Copy the local code to the container
COPY . .

# Build the Go application
RUN go build -o main .

# Expose port 8080 to the host machine
EXPOSE 8080

# Run the binary when the container starts
CMD ["./main"]
