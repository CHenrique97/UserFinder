# Start from a base image that includes the Go runtime
FROM golang:1.17-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module definition and download the dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go binary
RUN go build -o app

# Set the environment variables
ENV DSN=<insert DSN here>
ENV PORT=3030
ENV SECRET_KEY=<insert secret key here>

# Expose port 3030
EXPOSE 3030
EXPOSE 3306
# Start the Go microservice
CMD ["./app"]