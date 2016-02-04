package main

import (
	"bytes"
	"fmt"
	"github.com/MetalMatze/cubbyhole/server/cubbyhole"
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

	ResponseWelcome      = "!HELLO: Welcome to the Cubbyhole Server! Try 'help' for a list of commands" + ResponsePropmt
	ResponseHelp         = "!HELP:\nThe following commands are supported by this Cubbyhole:\n\nPUT <message>\t- Places a new message in the cubbyhole\nGET\t\t- Takes the message out of the cubbyhole and displays it\nLOOK\t\t- Displays the massage without taking it out of the cubbyhole\nDROP\t\t- Takes the message out of the cubbyhole without displaying it\nHELP\t\t- Displays this help message\nQUIT\t\t- Terminates the connection\n" + ResponsePropmt
	ResponseDrop         = "!DROP: ok" + ResponsePropmt
	ResponseGet          = "!GET: %s" + ResponsePropmt
	ResponseLook         = "!LOOK: %s" + ResponsePropmt
	ResponsePut          = "!PUT: ok" + ResponsePropmt
	ResponseQuit         = "!QUIT: ok"
	ResponseNotSupported = "!NOT SUPPORTED" + ResponsePropmt
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

	cubbyhole := cubbyhole.Cubbyhole{}

	app.Action = func(c *cli.Context) {
		listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
		defer listener.Close()
		if err != nil {
			log.Panic(err)
		}
		log.Printf("Listening on %s:%d", host, port)

		for {
			connection, err := listener.Accept()
			if err != nil {
				log.Panic(err)
			}
			log.Printf("Received %s -> %s", connection.RemoteAddr(), connection.LocalAddr())

			channel := make(chan string)

			go handleRequest(connection, channel, &cubbyhole)
			go sendData(connection, channel)

			channel <- ResponseWelcome
		}
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func handleRequest(connection net.Conn, channel chan string, cubbyhole *cubbyhole.Cubbyhole) {
	for {
		buffer := make([]byte, 1024)
		bufferLen, err := connection.Read(buffer)
		if err != nil {
			log.Panic(err)
		}

		request := string(buffer[:bufferLen])
		requestStrings := strings.Split(strings.TrimSpace(request), " ")

		switch strings.ToLower(requestStrings[0]) {
		case RequestPut:
			log.Println(connection.RemoteAddr(), RequestPut)
			cubbyhole.Put(strings.Join(requestStrings[1:len(requestStrings)], " "))
			channel <- ResponsePut
		case RequestGet:
			log.Println(connection.RemoteAddr(), RequestGet)
			if message := cubbyhole.Get(); message == "" {
				channel <- fmt.Sprintf(ResponseGet, ResponseNoMessage)
			} else {
				channel <- fmt.Sprintf(ResponseGet, message)
			}
		case RequestLook:
			log.Println(connection.RemoteAddr(), RequestLook)
			if message := cubbyhole.Look(); message == "" {
				channel <- fmt.Sprintf(ResponseLook, ResponseNoMessage)
			} else {
				channel <- fmt.Sprintf(ResponseLook, message)
			}
		case RequestDrop:
			log.Println(connection.RemoteAddr(), RequestDrop)
			cubbyhole.Drop()
			channel <- ResponseDrop
		case RequestHelp:
			log.Println(connection.RemoteAddr(), RequestHelp)
			channel <- ResponseHelp
		case RequestQuit:
			log.Println(connection.RemoteAddr(), RequestQuit)
			channel <- ResponseQuit
		default:
			channel <- ResponseNotSupported
		}
	}
}

func sendData(connection net.Conn, channel chan string) {
	for {
		response := <-channel
		_, err := io.Copy(connection, bytes.NewBufferString(response))
		if err != nil {
			log.Panic(err)
		}

		if response == ResponseQuit {
			connection.Close()
		}
	}
}
