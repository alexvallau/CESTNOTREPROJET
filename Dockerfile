# From golang image
FROM golang:latest
# Set the Current Working Directory inside the container
WORKDIR /app
# Copy the go.mod and go.sum files
COPY go.mod ./
COPY go.sum ./
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
# Copy the source code into the container
COPY . .
# Build the Go app
RUN go build -o main .
# Expose port 8080 to the outside world
EXPOSE 8080
# Command to run the executable
CMD ["./main"]
# Use the following command to build the Docker image
# docker build -t go-docker-example .

# Use the following command to run the Docker container
# docker run -p 8080:8080 go-docker-example