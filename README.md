# Cubbyhole server written in Go
A concurrent TCP server implement in Go serving clients speaking the [cubbyhole protocol](https://github.com/numbleroot/cubbyhole-server#protocol).

## Install with docker
Clone this repository to your machine.  

	docker build -t cubbyhole .
	docker run --publish 1337:1337 --name cubbyhole --rm cubbyhole

To stop the docker container run `docker stop cubbyhole`.

## Install by hand
Make sure you have [Go](https://golang.org) installed.  

	go get github.com/MetalMatze/cubbyhole

Now move to `cd $GOPATH/github.com/MetalMatze/cubbyhole`.

### run.sh
Use `./run.sh` to fetch dependencies and start the server.

### Without run.sh
First of all run `go get ./...` to fetch the few dependencies.  
Now you can start the server by running `go run main.go`.  
If you want to create a binary, do so by running `go build`.

## Other cubbyhole implementations

Of course I am not the only student in that class, so here is a probably incomplete list of implementations of my fellow students:

[cubbyhole](https://github.com/numbleroot/cubbyhole-server) by [numbleroot](https://github.com/numbleroot) - C  
[cubbyhole](https://github.com/KjellPirschel/cubbyhole) by [KjellPirschel](https://github.com/KjellPirschel) - C

## License
This project is [GPLv3](LICENSE) licensed.
