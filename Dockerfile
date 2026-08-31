# Use the official Golang image to create a build artifact.
# This is based on Alpine and sets the GOPATH to /go.
FROM golang:1.25-alpine as builder


# Create and change to the app directory.
WORKDIR /app

# Copy go.sum/go.mod and warm up the module cache.
COPY go.* ./

# Install git (use apk for Alpine, not apt-get)
RUN apk add --no-cache git


RUN go mod download

# Set the environment variable for Gin in release mode.
ENV GIN_MODE release

# Now copy the rest of the application's source code
COPY . .

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o server github.com/savannahghi/empower-clinical

FROM debian:bullseye-slim AS production

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

ENV USER=clinical
ENV WORKDIR=/app

RUN useradd --shell /bin/bash --no-create-home $USER

RUN mkdir -p $WORKDIR && chown -R $USER:$USER $WORKDIR
WORKDIR $WORKDIR

# Copy the Go binary to the production image from the builder stage.
COPY --from=builder /app/server /server

# Switch to non-root user
USER $USER

# Define the entrypoint for the container
CMD [ "/server" ]