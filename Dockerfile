FROM golang:1.15-alpine as builder

# Set environment variables to build the application binary for running on scratch image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o main cmd/stn-accounts/main.go

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/main .

# Build a small image. The overall virtual size of a running image is ~30mb
FROM scratch

COPY --from=builder /dist/main /

# Command to run
ENTRYPOINT ["/main"]