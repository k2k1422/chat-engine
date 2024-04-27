# Use an official Golang runtime as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

RUN apt update; \
        export DEBIAN_FRONTEND=noninteractive; \
        apt-get install -y tzdata; \
        dpkg-reconfigure --frontend noninteractive tzdata; \
        ln -fs /usr/share/zoneinfo/Asia/Kolkata /etc/localtime; \
        rm -rf /var/lib/apt/lists/*

# Copy the Go module files
COPY go.mod go.sum ./

# Download and install Go module dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Set CGO_ENABLED to 1 to enable CGO
ENV CGO_ENABLED=1

# Build the Go application
RUN go build -o myapp

# Expose any necessary ports
EXPOSE 8080

# Command to run the executable
CMD ["./myapp"]
