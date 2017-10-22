package tcpServer

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"436bin/a1/model"
	"encoding/json"
)

var rooms []*model.ChatRoom
var users []*model.User

const maxUsers int = 2

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
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("Error reading the room name from the request body")
		writer.WriteHeader(http.StatusInternalServerError)
	}
	var chatRoomName string
	json.Unmarshal(bodyBytes, &chatRoomName)
	room, err := getRoomForName(chatRoomName)
	if err != nil {
		fmt.Println(err.Error())
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	serializedLog, err := json.Marshal(&room.Log)
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

func listRoomsForUser(writer http.ResponseWriter, req *http.Request) {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("Error reading the user from the request body")
		writer.WriteHeader(http.StatusInternalServerError)
	}
	var userName string
	json.Unmarshal(bodyBytes, &userName)

	var roomsForUser []*model.ChatRoom
	for _, user := range users {
		if user.UserName == userName {
			roomsForUser = user.ChatRooms
		}
	}
	// if user not found, add the user
	if roomsForUser == nil {
		users = append(users, &model.User{UserName: userName})
	}
	serializedRooms, err := json.Marshal(&roomsForUser)
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

func joinRoom(writer http.ResponseWriter, req *http.Request) {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("Error reading the room from the request body")
		writer.WriteHeader(http.StatusInternalServerError)
	}
	var joinRequest *model.JoinChatRequest
	json.Unmarshal(bodyBytes, &joinRequest)
	for _, room := range rooms {
		if room.Name != joinRequest.RoomName {
			continue
		}
		// Check, lock, check
		if len(room.Users) >= room.MaxUsers {
			writer.WriteHeader(http.StatusForbidden)
			return
		}
		room.Mux.Lock()
		if len(room.Users) >= room.MaxUsers {
			room.Mux.Unlock()
			writer.WriteHeader(http.StatusForbidden)
			return
		}
		room.Users = append(room.Users, joinRequest.User)
		room.Mux.Unlock()
	}
}

// *************

func logMessage(m *model.Message) {
	fmt.Printf("%s: %s\n", string(m.UserName), string(m.Body))
	room, err := getRoomForName(m.ChatRoomName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	room.Log = append(room.Log, m)
}

func getRoomForName(chatRoomName string) (*model.ChatRoom, error) {
	for _, room := range rooms {
		if room.Name == chatRoomName {
			return room, nil
		}
	}
	return nil, errors.New("Failed to find chat room")
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
	http.HandleFunc("/chatrooms/join", joinRoom)
	http.HandleFunc("/chatrooms/forUser", listRoomsForUser)
	// http.HandleFunc("/chatrooms/leave", todo)

	start()
}
