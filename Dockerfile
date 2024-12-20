FROM golang:1.23-alpine as builder

LABEL maintainer="backend-ritlab"
LABEL description="Golang Service for handling Text to Speech"
RUN apk add --no-cache make \
     build-base tzdata musl-dev gcc g++ libc6-compat

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

ENV SERVICE_PORT=7075

# Create and set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY docker/app .

RUN echo $(ls -al)

RUN go env
RUN go build -tags musl -x -v -o tts cmd/main.go


FROM ubuntu:22.04

# Set environment variables
ENV DEBIAN_FRONTEND=noninteractive
ARG PROFILE
ENV APP_ENV=${PROFILE}

# Install necessary runtime dependencies
RUN apt update && apt install -y --no-install-recommends \
    mpg123 && \
    rm -rf /var/lib/apt/lists/*

EXPOSE ${SERVICE_PORT}
WORKDIR /app

COPY --from=builder /app/tts /app/
COPY --from=builder /app/audio/ /app/audio

CMD ["./tts"]