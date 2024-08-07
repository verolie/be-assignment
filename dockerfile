# Use official Golang image as base
FROM golang:1.18

# Install Node.js 16 and npm
RUN apt-get update && \
    apt-get install -y curl && \
    curl -fsSL https://deb.nodesource.com/setup_16.x | bash - && \
    apt-get install -y nodejs

# Set the working directory inside the container
WORKDIR /app

# Copy Go modules and install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy all Go source files and directories
COPY . .

# Install Prisma CLI
RUN npm install -g prisma

# Generate Prisma Client (Go)
RUN go get github.com/steebchen/prisma-client-go

# Build the Go application
RUN go build -o main .

# Command to run the application
CMD ["./main"]