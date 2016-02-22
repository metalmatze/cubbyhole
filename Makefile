all: build

deps:
	go get -u github.com/Masterminds/glide
	glide install

build: 
	go build -o server

clean:
	rm -r server
