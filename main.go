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

			client.Outgoing <- ResponseWelcome
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
		case RequestPut:
			log.Println(client.Connection.RemoteAddr(), RequestPut)
			cubbyhole.Put(strings.Join(requestStrings[1:], " "))
			client.Outgoing <- ResponsePut
		case RequestGet:
			log.Println(client.Connection.RemoteAddr(), RequestGet)
			if message := cubbyhole.Get(); message == "" {
				client.Outgoing <- fmt.Sprintf(ResponseGet, ResponseNoMessage)
			} else {
				client.Outgoing <- fmt.Sprintf(ResponseGet, message)
			}
		case RequestLook:
			log.Println(client.Connection.RemoteAddr(), RequestLook)
			if message := cubbyhole.Look(); message == "" {
				client.Outgoing <- fmt.Sprintf(ResponseLook, ResponseNoMessage)
			} else {
				client.Outgoing <- fmt.Sprintf(ResponseLook, message)
			}
		case RequestDrop:
			log.Println(client.Connection.RemoteAddr(), RequestDrop)
			cubbyhole.Drop()
			client.Outgoing <- ResponseDrop
		case RequestHelp:
			log.Println(client.Connection.RemoteAddr(), RequestHelp)
			client.Outgoing <- ResponseHelp
		case RequestQuit:
			log.Println(client.Connection.RemoteAddr(), RequestQuit)
			client.Outgoing <- ResponseQuit
		default:
			client.Outgoing <- ResponseNotSupported
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

		if response == ResponseQuit {
			client.Connection.Close()
			client.Close()
		}
	}
}
