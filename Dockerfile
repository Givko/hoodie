# Build stage
FROM golang:1.23.4-bullseye AS builder

WORKDIR /usr/src/app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify && go mod tidy

# Copy source code
COPY . .

# Build with optimizations
#RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /usr/local/bin/app .
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/app .

# Final stage
FROM gcr.io/distroless/static-debian11

# Add labels
LABEL maintainer="milevjivko@gmail.com"
LABEL version="0.0.1"
LABEL description="Hoodie a simple friendly chat server"

COPY --from=builder /usr/local/bin/app /app

# Use non-root user; verify UID and GID if necessary
USER 65532:65532

EXPOSE 8080

ENTRYPOINT ["/app"]
