# syntax=docker/dockerfile:1
FROM golang:latest

# Set destination for COPY
WORKDIR /app

# Copy files
COPY . ./

# Download dependencies
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build main.go

# Run
CMD ["./main"]