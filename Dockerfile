# Build stage
FROM golang:1.22 as builder

WORKDIR /app

# Copy the Go Modules manifests
COPY go.mod go.sum ./
# Download Go modules
RUN go mod download

# Copy the go source code
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Final stage
FROM scratch

# Set the working directory in the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/app .

EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["./app"]
