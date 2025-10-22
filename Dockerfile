# Use the official lightweight Go image based on Alpine Linux
FROM golang:1.24-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy all project files from your local machine into the container
COPY . .

# Download Go module dependencies (like go.mod and go.sum)
RUN go mod tidy

# Expose port 9000 so it can be accessed from outside the container
EXPOSE 9000

# Run the application directly from source code (no build needed)
CMD ["go", "run", "main.go"]
