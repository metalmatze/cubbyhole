package main

import (
	"bytes"
	"fmt"
	"github.com/MetalMatze/cubbyhole/cubbyhole"
	"github.com/codegangsta/cli"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

const (
	requestPut  = "put"
	requestGet  = "get"
	requestLook = "look"
	requestDrop = "drop"
	requestHelp = "help"
	requestQuit = "quit"

	responseWelcome      = "!HELLO: Welcome to the Cubbyhole Server! Try 'help' for a list of commands" + responsePropmt
	responseHelp         = "!HELP:\nThe following commands are supported by this Cubbyhole:\n\nPUT <message>\t- Places a new message in the cubbyhole\nGET\t\t- Takes the message out of the cubbyhole and displays it\nLOOK\t\t- Displays the massage without taking it out of the cubbyhole\nDROP\t\t- Takes the message out of the cubbyhole without displaying it\nHELP\t\t- Displays this help message\nQUIT\t\t- Terminates the connection\n" + responsePropmt
	responseDrop         = "!DROP: ok" + responsePropmt
	responseGet          = "!GET: %s" + responsePropmt
	responseLook         = "!LOOK: %s" + responsePropmt
	responsePut          = "!PUT: ok" + responsePropmt
	responseQuit         = "!QUIT: ok"
	responseNotSupported = "!NOT SUPPORTED" + responsePropmt
	responseNoMessage    = "<no message stored>"
	responsePropmt       = "\n> "
)

var (
	host = ""
	port = 1337
)

// Client is a wrapper for its net.Conn and channels.
type Client struct {
	Connection net.Conn
	Connected  chan bool
	Incoming   chan string
	Outgoing   chan string
}

// ReadString reads from the client's connection into a buffer and returns the buffer's content as string.
func (c *Client) ReadString() (string, error) {
	buffer := make([]byte, 1024)
	bytesRead, err := c.Connection.Read(buffer)
	if err != nil {
		c.Close()
		log.Fatal(err)
		return "", err
	}
	request := string(buffer[:bytesRead])
	return request, nil
}

// Close shutsdown the client's connection and its channels.
func (c *Client) Close() {
	log.Printf("Closing connection to %s", c.Connection.RemoteAddr())
	c.Connected <- false
	close(c.Incoming)
	close(c.Outgoing)
	close(c.Connected)
	if err := c.Connection.Close(); err != nil {
		log.Fatal(err)
	}
}

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
			log.Printf("Listenting %s -> %s", connection.RemoteAddr(), connection.LocalAddr())

			client := Client{
				Connection: connection,
				Connected:  make(chan bool),
				Incoming:   make(chan string),
				Outgoing:   make(chan string),
			}

			go handleRequest(&client, &cubbyhole)
			go sendData(&client)

			client.Outgoing <- responseWelcome
		}
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func handleRequest(client *Client, cubbyhole *cubbyhole.Cubbyhole) {
	for {
		request, err := client.ReadString()
		if err != nil {
			log.Println(err)
			return
		}

		requestStrings := strings.Split(strings.TrimSpace(request), " ")
		switch strings.ToLower(requestStrings[0]) {
		case requestPut:
			log.Println(client.Connection.RemoteAddr(), requestPut)
			cubbyhole.Put(strings.Join(requestStrings[1:], " "))
			client.Outgoing <- responsePut
		case requestGet:
			log.Println(client.Connection.RemoteAddr(), requestGet)
			if message := cubbyhole.Get(); message == "" {
				client.Outgoing <- fmt.Sprintf(responseGet, responseNoMessage)
			} else {
				client.Outgoing <- fmt.Sprintf(responseGet, message)
			}
		case requestLook:
			log.Println(client.Connection.RemoteAddr(), requestLook)
			if message := cubbyhole.Look(); message == "" {
				client.Outgoing <- fmt.Sprintf(responseLook, responseNoMessage)
			} else {
				client.Outgoing <- fmt.Sprintf(responseLook, message)
			}
		case requestDrop:
			log.Println(client.Connection.RemoteAddr(), requestDrop)
			cubbyhole.Drop()
			client.Outgoing <- responseDrop
		case requestHelp:
			log.Println(client.Connection.RemoteAddr(), requestHelp)
			client.Outgoing <- responseHelp
		case requestQuit:
			log.Println(client.Connection.RemoteAddr(), requestQuit)
			client.Outgoing <- responseQuit
		default:
			client.Outgoing <- responseNotSupported
		}
	}
}

func sendData(client *Client) {
	for {
		response := <-client.Outgoing
		_, err := io.Copy(client.Connection, bytes.NewBufferString(response))
		if err != nil {
			log.Panic(err)
		}

		if response == responseQuit {
			client.Connection.Close()
			client.Close()
		}
	}
}
