# Cubbyhole server written in Go
A concurrent TCP server implement in Go serving clients speaking the [cubbyhole protocol](https://github.com/numbleroot/cubbyhole-server#protocol).

## Installation

	go get -u github.com/MetalMatze/cubbyhole

Now move to `cd $GOPATH/src/github.com/MetalMatze/cubbyhole`.

Simply run `make` which will fetch the few dependencies, build the binary & test it afterwards.

Start the server by running `./server`

### Installation with Docker
Clone this repository to your machine.  

	docker build -t cubbyhole .
	docker run --publish 1337:1337 --name cubbyhole --rm cubbyhole

To stop the docker container run `docker stop cubbyhole`.

## Other cubbyhole implementations

Of course I am not the only student in that class, so here is a probably incomplete list of implementations of my fellow students:

[cubbyhole](https://github.com/numbleroot/cubbyhole-server) by [numbleroot](https://github.com/numbleroot) - C  
[cubbyhole](https://github.com/KjellPirschel/cubbyhole) by [KjellPirschel](https://github.com/KjellPirschel) - C

## License
This project is [GPLv3](LICENSE) licensed.
