# Using official golang image as base image
FROM golang:1.23.2-alpine

# Creating work directory   
WORKDIR /app

# # Copying go.mod and go.sum files first, then running go mod download 
# COPY go.mod go.sum ./ 
# RUN go mod download

# Copying all files from the current directory to the container
COPY . .

# Downloading dependencies
RUN go mod tidy
RUN go build -o cashin-cashout ./cmd/cashin-cashout/main.go

# Command to run the executable
CMD ["./cashin-cashout"]