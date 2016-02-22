all: build

PACKAGES = $(shell go list ./... | grep -v /vendor/)

deps:
	go get -u github.com/Masterminds/glide
	glide install

build: 
	go build -o server

clean:
	rm -r server

test:
	go test -cover $(PACKAGES)
