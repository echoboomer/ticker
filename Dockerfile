FROM golang:alpine AS builder

ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

WORKDIR /build

# Download go dependencies
COPY go.mod .
COPY go.sum .
COPY static/ .
RUN go mod download

# Copy into the container
COPY . .

# Build the application
RUN go build -o ticker .

# Build final image using nothing but the binary
FROM alpine:3.16.2

COPY --from=builder /build/ticker /
COPY --from=builder /build/static /static
COPY --from=builder /build/docs /docs

EXPOSE 8080

# Command to run
CMD ["/ticker"]