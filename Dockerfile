# Build Stage
FROM golang:1.24-alpine AS builder

# Build dependencies skipped (network issue workaround)
# RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download skipped because we use local vendor
# RUN go mod download

# Copy the source code (including vendor folder)
COPY . .

# Build the application
# Use -mod=vendor to build using local dependencies (offline mode)
RUN CGO_ENABLED=0 go build -mod=vendor -ldflags="-w -s" -o learn main.go

# Final Stage
FROM alpine:latest

# Install necessary runtime dependencies
# Note: if network fails here too, we might need to skip this or use offline method, 
# but certificates are usually needed. Let's hope simple apk works or is cached.
RUN apk --no-cache add ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/learn .

# Expose the application port
EXPOSE 8080

# Create a non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Set the entrypoint to run the application
ENTRYPOINT ["./learn"]
CMD ["serve"]