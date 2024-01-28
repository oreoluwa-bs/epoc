# Use a smaller base image
FROM golang:alpine as builder

# Set the working directory inside the container
WORKDIR /app

# Install necessary packages for building with CGO
RUN apk add --no-cache gcc musl-dev

# Copy only the necessary files
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the local package files to the container's workspace
COPY . .

# Build the Go application with CGO enabled
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o server .

# Use a minimal Alpine image as the final base
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy only the necessary files from the builder stage
COPY --from=builder /app/server .

RUN mkdir data

# Expose the port on which the server will run
EXPOSE 8000

# Command to run the executable
CMD ["./server"]
