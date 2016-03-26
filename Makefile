PACKAGES = $(shell go list ./... | grep -v /vendor/)

all: deps build test

deps:
	go get -u github.com/govend/govend
	govend -v

build: 
	go build -o server

clean:
	rm -r server

test:
	go test -cover $(PACKAGES)
