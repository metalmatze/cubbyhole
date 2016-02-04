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
	RequestPut  = "put"
	RequestGet  = "get"
	RequestLook = "look"
	RequestDrop = "drop"
	RequestHelp = "help"
	RequestQuit = "quit"

	ResponseWelcome      = "!HELLO: Welcome to the Cubbyhole Server! Try 'help' for a list of commands"
	ResponseHelp         = "!HELP:\nThe following commands are supported by this Cubbyhole:\n\nPUT <message>\t- Places a new message in the cubbyhole\nGET\t\t- Takes the message out of the cubbyhole and displays it\nLOOK\t\t- Displays the massage without taking it out of the cubbyhole\nDROP\t\t- Takes the message out of the cubbyhole without displaying it\nHELP\t\t- Displays this help message\nQUIT\t\t- Terminates the connection\n"
	ResponseDrop         = "!DROP: ok"
	ResponseGet          = "!GET: "
	ResponseLook         = "!LOOK: "
	ResponsePut          = "!PUT: ok"
	ResponseQuit         = "!QUIT: ok"
	ResponseNotSupported = "!NOT SUPPORTED"
	ResponseNoMessage    = "<no message stored>"
	ResponsePropmt       = "\n> "
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

			go handleRequest(connection, channel)
			go sendData(connection, channel)

			channel <- ResponseWelcome + ResponsePropmt
		}
	}

	app.Run(os.Args)
}

func handleRequest(connection net.Conn, channel chan string) {
	for {
		buffer := make([]byte, 1024)
		len, err := connection.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}

		request := string(buffer[:len])
		requestStrings := strings.Split(strings.TrimSpace(request), " ")

		switch strings.ToLower(requestStrings[0]) {
		case RequestPut:
			channel <- "putting"
		case RequestGet:
			channel <- "getting"
		case RequestLook:
			channel <- "looking"
		case RequestDrop:
			channel <- "dropping"
		case RequestHelp:
			log.Println(connection.RemoteAddr(), "help")
			channel <- ResponseHelp + ResponsePropmt
		case RequestQuit:
			channel <- ResponseQuit
		default:
			channel <- ResponseNotSupported + ResponsePropmt
		}
	}
}

func sendData(connection net.Conn, channel chan string) {
	for {
		response := <-channel
		io.Copy(connection, bytes.NewBufferString(response))
	}
}
