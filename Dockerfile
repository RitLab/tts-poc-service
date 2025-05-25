# Stage 1: Builder
FROM golang:1.23-bullseye AS builder

LABEL maintainer="backend-ritlab"
LABEL description="Golang Service for handling Text to Speech"

# Install dependencies for ALSA
RUN apt-get update && apt-get install -y \
    libasound2-dev \
    build-essential \
    pkg-config \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Set environment variables for CGO
ENV CGO_ENABLED=1 \
    GO111MODULE=on

# Set the working directory
WORKDIR /app

# Copy the Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the application source code
COPY . .

# Build the Go application
RUN go build -tags musl -x -v -o tts cmd/main.go

# Stage 2: Final Image
FROM debian:bullseye-slim

ARG PROFILE
ARG CONFIG
ENV APP_ENV=${PROFILE}
ENV CONFIG_FILE=${CONFIG}

# Install runtime dependencies for ALSA
RUN apt-get update && apt-get install -y \
    libasound2 \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

EXPOSE 7075

# Set the working directory
WORKDIR /app

# Copy the built executable from the builder stage
COPY --from=builder /app/tts /app/
COPY --from=builder /app/audio/ /app/audio
COPY --from=builder /app/pdf/ /app/pdf
COPY --from=builder /app/public/ /app/public
COPY --from=builder /app/api/ /app/api

RUN chmod -R 777 /app/audio
RUN chmod -R 777 /app/pdf

# Set the entrypoint
CMD ["./tts"]