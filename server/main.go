package main

import (
	"bytes"
	"fmt"
	"github.com/codegangsta/cli"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

const (
	RESPONSE_WELCOME       = "!HELLO: Welcome to the Cubbyhole Server! Try 'help' for a list of commands"
	RESPONSE_HELP          = "!HELP:\nThe following commands are supported by this Cubbyhole:\n\nPUT <message>\t- Places a new message in the cubbyhole\nGET\t\t- Takes the message out of the cubbyhole and displays it\nLOOK\t\t- Displays the massage without taking it out of the cubbyhole\nDROP\t\t- Takes the message out of the cubbyhole without displaying it\nHELP\t\t- Displays this help message\nQUIT\t\t- Terminates the connection\n"
	RESPONSE_DROP          = "!DROP: ok"
	RESPONSE_GET           = "!GET: "
	RESPONSE_LOOK          = "!LOOK: "
	RESPONSE_PUT           = "!PUT: ok"
	RESPONSE_QUIT          = "!QUIT: ok"
	RESPONSE_NOT_SUPPORTED = "!NOT SUPPORTED"
	RESPONSE_NO_MESSAGE    = "<no message stored>"
	RESPONSE_PROPMT        = "\n> "
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

			channel <- RESPONSE_WELCOME + RESPONSE_PROPMT
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
		requestStrings := strings.Split(strings.TrimSpace(request), " ")

		switch strings.ToLower(requestStrings[0]) {
		case "put":
			channel <- "putting"
		case "get":
			channel <- "getting"
		case "look":
			channel <- "looking"
		case "drop":
			channel <- "dropping"
		case "help":
			channel <- RESPONSE_HELP + RESPONSE_PROPMT
		case "quit":
			connection.Close()
		default:
			channel <- RESPONSE_NOT_SUPPORTED + RESPONSE_PROPMT
		}
	}
}

func sendData(connection net.Conn, channel chan string) {
	for {
		response := <-channel
		io.Copy(connection, bytes.NewBufferString(response))
	}
}
