package tcpClient

import (
	"436bin/a1/model"
	"encoding/gob"

	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// ReadUserName retrieves the desired username
type ReadUserName func(client *Client) (string, error)

// Client contains our username, and helper methods
type Client struct {
	UserName     string
	readUserName ReadUserName
	Conn         *net.Conn
}

func readUserName(client *Client) (string, error) {
	reader := createBufioReader()
	fmt.Print("Please enter your desired user name: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return text, nil
}

// CreateUser constructs the user with an associated UserName
func (client *Client) CreateUser() error {
	text, err := client.readUserName(client)
	if err != nil {
		return err
	}
	text = strings.TrimSpace(text)
	client.UserName = text
	return nil
}

// CreateReader this is split out for testing purposes
func createBufioReader() *bufio.Reader {
	return bufio.NewReader(os.Stdin)
}

func constructMessage(userName string, body string) *model.Message {
	return &model.Message{Sender: userName, Body: strings.TrimSpace(body)}
}

func printMessage(m *model.Message) {
	fmt.Printf("%s: %s\n", m.Sender, m.Body)
}

func sendMessage(client *Client, text string) {
	// Format the message for serialization
	m := constructMessage(client.UserName, text)
	fmt.Printf("Message to send: %s\n", m) // todo: remove this
	// Use gob lib to encode the data
	enc := gob.NewEncoder(*client.Conn) // to write
	enc.Encode(&m)
}

func receiveMessage(client *Client) {
	message := &model.Message{}
	dec := gob.NewDecoder(*client.Conn)
	err := dec.Decode(message)
	if err != nil {
		fmt.Println("Decoding response from server failed.")
		return
	}

	printMessage(message)
}

func readMessageFromUser(client *Client) (string, error) {
	fmt.Printf("%s: ", client.UserName)
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return text, nil
}

// Create makes a new tcp client and waits to send a message to the target server.
func Create() {
	fmt.Println("Creating client...")
	// Create the client
	client := &Client{readUserName: readUserName}

	// Before dialing, we set up the username
	for {
		err := client.CreateUser()
		if err != nil {
			fmt.Println("Failed to create user, please try again.")
			continue
		}
		break
	}

	c, err := net.Dial("tcp", "127.0.0.1:8081")
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Client dialing failed. Exiting.")
		os.Exit(1)
	}
	client.Conn = &c

	// Receive the log, and print it
	dec := gob.NewDecoder(c) // to read
	var log []*model.Message
	err = nil
	err = dec.Decode(&log)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}

	for _, logMessage := range log {
		printMessage(logMessage)
	}

	for {
		text, err := readMessageFromUser(client)
		if err != nil {
			fmt.Println("Failed to accept message input. Please try again.")
			continue
		}

		sendMessage(client, text)

		receiveMessage(client)
	}
}
