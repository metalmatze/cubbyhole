# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1-alpine

WORKDIR /go/src/github.com/MetalMatze/cubbyhole

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/MetalMatze/cubbyhole

RUN apk update && apk upgrade
RUN apk add git

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go get -v ./...
RUN go install github.com/MetalMatze/cubbyhole

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/cubbyhole

# Document that the service listens on port 1337.
EXPOSE 1337
