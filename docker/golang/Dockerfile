# Use the latest Go image as the builder
FROM golang:latest AS builder

# Install Delve (dlv)
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Create the application directory
RUN mkdir /app

# Add the application source code to the container
ADD . /app

# Set the working directory
WORKDIR /app

# Build the application with debugging flags
RUN CGO_ENABLED=0 GOOS=linux go build -gcflags="all=-N -l" -o main ./...

# Use the latest Alpine image for the production stage
FROM alpine:latest

# Copy the built application and Delve from the builder stage
COPY --from=builder /app .
COPY --from=builder /go/bin/dlv /usr/local/bin/dlv

# Expose the necessary ports
EXPOSE 8000 40000

# Set environment variables
ENV PORT=8000

# Command to run the application with Delve
CMD ["/usr/local/bin/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "./main"]