# Use the official lightweight Go image based on Alpine Linux
FROM golang:1.24-alpine

# Install necessary tools
RUN apk add --no-cache git

# Install swag CLI globally
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Set the working directory inside the container
WORKDIR /app

# Copy all project files from your local machine into the container
COPY . .

# Download Go module dependencies (like go.mod and go.sum)
RUN go mod tidy

# Generate Swagger
RUN swag init

# Expose port 9000 so it can be accessed from outside the container
EXPOSE 9000

# Add Go bin to PATH so swag is accessible
ENV PATH="/go/bin:${PATH}"

# Run the application directly from source code (no build needed)
CMD ["go", "run", "main.go"]
