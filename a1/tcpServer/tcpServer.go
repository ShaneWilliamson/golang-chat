package tcpServer

import (
	"errors"
	"fmt"
	"net"
	"os"

	"436bin/a1/model"
	"encoding/gob"
)

func acceptNewConn(link net.Listener) (net.Conn, error) {
	return link.Accept()
}

func readFromConnection(c net.Conn) (model.Message, error) {
	var retries = 10
	for i := 0; i < retries; i++ {
		// Read in the message from the client
		dec := gob.NewDecoder(c)
		var message model.Message
		err := dec.Decode(message)
		// message, err := bufio.NewReader(c).ReadString('\n') //TODO: cleanup
		if err == nil {
			return message, nil
		}
	}
	// return an error to handle reopening of connection
	return model.Message{Sender: "", Body: ""}, errors.New("Failed to read from connection.")
}

func start(link net.Listener) {
	c, err := acceptNewConn(link)
	if err != nil {
		fmt.Println("Error connecting on port 8081")
		os.Exit(2)
	}
	// var log []model.Message
	for {
		// will listen for message to process ending in newline (\n)
		m, err := readFromConnection(c)
		if err != nil {
			// Reading from our connection has failed, accept a new one
			if err.Error() == "Failed to read from connection." {
				// This is an acceptable error case, here's where we accept a new connection
				// we will accept a new connection if we've failed with retries once, otherwise exit for now
				fmt.Println("Connection to client has been closed. Listening for new connection.")
				c, err = acceptNewConn(link)
				if err != nil {
					// Creation of new connection failed
					fmt.Println("Failed to accept a new connection. Exiting.")
					os.Exit(3)
				}
			}
			// we've created a new connection, now continue loop to attempt reading from a client
			continue
		}

		fmt.Printf("%s: %s", string(m.Sender), string(m.Body))
		// TODO log it and then send response to client.

		// // output message received
		// fmt.Print("Message Received:", string(m))
		// // add the message to the log
		// log = append(log, model.Message{"foo", m})
		// // sample process for string received
		// newmessage := strings.ToUpper(m)
		// // send new string back to client
		// c.Write([]byte(newmessage + "\n"))
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
