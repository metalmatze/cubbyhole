package main

import (
	"bytes"
	"fmt"
	"github.com/codegangsta/cli"
	"io"
	"log"
	"net"
	"os"
)

var (
	host string = ""
	port int    = 1337
)

func main() {
	app := cli.NewApp()
	app.Name = "cubbyhole Server"
	app.HideVersion = true
	app.Usage = "A cubbyhole server in go"

	app.Action = func(c *cli.Context) {
		listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
		defer listener.Close()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Listening on %s:%d", host, port)

		for {
			connection, err := listener.Accept()
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Received %s -> %s", connection.RemoteAddr(), connection.LocalAddr())

			channel := make(chan string)
			errChannel := make(chan error)

			go handleRequest(connection, channel, errChannel)
			go sendData(connection, channel)
		}
	}

	app.Run(os.Args)
}

func handleRequest(connection net.Conn, channel chan string, errChannel chan error) {
	for {
		buffer := make([]byte, 1024)
		len, err := connection.Read(buffer)
		if err != nil {
			errChannel <- err
		}

		request := string(buffer[:len])
		channel <- request
	}
}

func sendData(connection net.Conn, channel chan string) {
	for {
		response := <-channel
		io.Copy(connection, bytes.NewBufferString(response))
	}
}
