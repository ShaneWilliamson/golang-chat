package tcpServer

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"

	"436bin/a1/model"
	"encoding/gob"
)

var log []*model.Message

func sendLogToConn(c *net.Conn) {
	enc := gob.NewEncoder(*c) // to write
	err := enc.Encode(log)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func acceptNewConn(link net.Listener) net.Conn {
	c, err := link.Accept()
	if err != nil {
		// Creation of new connection failed
		fmt.Println("Failed to accept a new connection. Exiting.")
		os.Exit(3)
	}
	// Send log to new connection
	sendLogToConn(&c)

	return c
}

func readFromConnection(c net.Conn) (*model.Message, error) {
	var retries = 10
	for i := 0; i < retries; i++ {
		// Read in the message from the client
		dec := gob.NewDecoder(c)
		message := &model.Message{}
		err := dec.Decode(message)
		if err == io.EOF {
			return nil, errors.New("Connection closed")
		}
		if err != nil {
			// We should still reply with an error message
			fmt.Println(err.Error())
			enc := gob.NewEncoder(c)
			// Todo: this won't cut it, needs to be a pointer to some object. Add err field to message
			enc.Encode("ERROR: Please try again in a few moments")
		}

		if err == nil {
			return message, nil
		}
	}
	// return an error to handle reopening of connection
	return &model.Message{}, errors.New("Failed to read from connection")
}

func logMessage(m *model.Message) {
	fmt.Printf("%s: %s\n", string(m.Sender), string(m.Body))
	log = append(log, m)
}

func respondToClient(c *net.Conn, m *model.Message) {
	enc := gob.NewEncoder(*c) // to write
	enc.Encode(m)
}

func start(link net.Listener) {
	fmt.Println("Starting server...")
	c := acceptNewConn(link)
	for {
		// will listen for message to process ending in newline (\n)
		m, err := readFromConnection(c)
		if err != nil {
			// Reading from our connection has failed, accept a new one
			if err.Error() == "Connection closed" {
				// This is an acceptable error case, here's where we accept a new connection
				// we will accept a new connection if we've failed with retries once, otherwise exit for now
				fmt.Println("Connection to client has been closed. Listening for new connection.")
				c = acceptNewConn(link)
			}
			// we've created a new connection, now continue loop to attempt reading from a client
			continue
		}
		logMessage(m)
		respondToClient(&c, m)
	}
}

// Create makes a new tcp server and listens for incoming requests
func Create() {
	// create the server
	fmt.Println("Creating Server...")
	link, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Println("Error attempting to listen on port 8081. Exiting.")
		os.Exit(1)
	}
	start(link)
}
