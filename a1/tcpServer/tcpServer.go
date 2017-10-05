package tcpServer

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"

	"436bin/a1/model"
	"encoding/gob"
	"encoding/json"
)

var log []*model.Message // This will be removed when we implement rooms

const threadCount int = 10

func sendLogToConn(c *net.Conn) {
	// TODO convert and serialize for http request
	enc := gob.NewEncoder(*c) // to write
	err := enc.Encode(log)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// TODO delete this after transition to http requests
func acceptNewConn(link net.Listener) net.Conn {
	// Wait for the next call, and returns a generic connection
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

func receiveMessage(writer http.ResponseWriter, req *http.Request) {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("Error reading the message from the request body")
		writer.WriteHeader(http.StatusInternalServerError)
	}
	var message *model.Message
	json.Unmarshal(bodyBytes, &message)
	logMessage(message)
}

func getLog(writer http.ResponseWriter, req *http.Request) {
	serializedLog, err := json.Marshal(&log)
	if err != nil {
		fmt.Println("Marshalling the log has failed.")
		writer.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Println(string(serializedLog))
	writer.Write(serializedLog)
}

func logMessage(m *model.Message) {
	fmt.Printf("%s: %s\n", string(m.Sender), string(m.Body))
	log = append(log, m)
}

func respondToClient(c *net.Conn, m *model.Message) {
	enc := gob.NewEncoder(*c) // to write
	enc.Encode(m)
}

func start() {
	fmt.Println("Starting server...")
	// Create the HTTP server
	fmt.Println((http.ListenAndServe(":8081", nil).Error()))
}

// Create makes a new tcp server and listens for incoming requests
func Create() {
	// create the server
	fmt.Println("Creating Server...")

	// Register our HTTP routes
	http.HandleFunc("/message", receiveMessage)
	http.HandleFunc("/log", getLog)
	// http.HandleFunc("/chatrooms/list", todo)
	// http.HandleFunc("/chatrooms/join", todo)
	// http.HandleFunc("/chatrooms/leave", todo)

	start()
}
