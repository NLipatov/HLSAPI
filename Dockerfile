# syntax=docker/dockerfile:1
FROM golang:latest

# Set the app root folder
ENV APP_ROOT=/app

# Set destination for COPY
WORKDIR /app

# Install FFMPEG
RUN apt update && apt install ffmpeg -y

# Copy files
COPY . ./

# Download dependencies
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build main.go

# Run
CMD ["./main"]