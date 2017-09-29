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

func constructMessage(userName string, body string) model.Message {
	return model.Message{Sender: userName, Body: strings.TrimSpace(body)}
}

// Create makes a new tcp client and waits to send a message to the target server.
func Create() {
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
	for {
		fmt.Print("sending message: ")
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Print("Error while reading in message. Exiting.")
			os.Exit(1)
		}
		// Format the message for serialization
		m := constructMessage(client.UserName, text)
		// Use gob lib to encode the data
		enc := gob.NewEncoder(c) // to write
		dec := gob.NewDecoder(c) // to read
		enc.Encode(m)
		// serializedMessage := model.SerializeMessage(m)
		if err != nil {
			fmt.Println("Message failed to serialize, please try again.")
			continue
		}

		var message model.Message
		err = dec.Decode(message)
		if err != nil {
			fmt.Println("Decoding response from server failed.")
			continue
		}

		// ** OLD METHOD **
		// Send the serialized message to the server
		// fmt.Fprintf(c, serializedMessage+"\n")
		// message, err := bufio.NewReader(c).ReadString('\n')
		// if err != nil {
		// 	fmt.Print("Failed to read response from server. Exiting.")
		// 	os.Exit(2)
		// }
		// **

		fmt.Println("Response from server:")
		fmt.Printf("Message: {%s, %s}\n", message.Sender, message.Body)
	}
}
