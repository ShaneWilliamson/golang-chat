package tcpServer

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"436bin/a1/model"
	"encoding/json"
)

var log []*model.Message // This will be removed when we implement rooms
var rooms []*model.ChatRoom

const maxUsers int = 10

func receiveMessage(writer http.ResponseWriter, req *http.Request) {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("Error reading the message from the request body")
		writer.WriteHeader(http.StatusInternalServerError)
	}
	var message *model.Message
	json.Unmarshal(bodyBytes, &message)
	logMessage(message)
	broadcastMessage(message)
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

func listRooms(writer http.ResponseWriter, req *http.Request) {
	serializedRooms, err := json.Marshal(&rooms)
	if err != nil {
		fmt.Println("Marshalling the rooms has failed.")
		writer.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Println(string(serializedRooms))
	writer.Write(serializedRooms)
}

func createRoom(writer http.ResponseWriter, req *http.Request) {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("Error reading the room from the request body")
		writer.WriteHeader(http.StatusInternalServerError)
	}
	var chatRoomName string
	json.Unmarshal(bodyBytes, &chatRoomName)
	rooms = append(rooms, &model.ChatRoom{Users: nil, Name: chatRoomName, MaxUsers: maxUsers})
}

// *************

func logMessage(m *model.Message) {
	fmt.Printf("%s: %s\n", string(m.UserName), string(m.Body))
	log = append(log, m)
}

func broadcastMessage(message *model.Message) {
	client := &http.Client{}
	// TODO broadcast to all users in the target room
	// Format the message for serialization
	messageBuffer := model.ConvertMessageToBuffer(message)
	client.Post("http://localhost:9081/message", "application/json; charset=utf-8", messageBuffer)
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
	http.HandleFunc("/chatrooms/list", listRooms)
	http.HandleFunc("/chatrooms/create", createRoom)
	// http.HandleFunc("/chatrooms/join", todo)
	// http.HandleFunc("/chatrooms/leave", todo)

	start()
}
