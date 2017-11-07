package tcpServer

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"436bin/a1/config"
	"436bin/a1/model"
	"encoding/json"
)

var rooms []*model.ChatRoom
var users []*model.User

const maxUsers int = 10

func receiveMessage(writer http.ResponseWriter, req *http.Request) {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("Error reading the message from the request body")
		writer.WriteHeader(http.StatusInternalServerError)
	}
	var message *model.Message
	json.Unmarshal(bodyBytes, &message)
	room, err := getRoomForName(message.ChatRoomName)
	if err != nil {
		fmt.Println(err.Error())
		writer.WriteHeader(http.StatusForbidden)
		return
	}
	if room.GetUser(message.UserName) != nil {
		go logMessage(message)
		go broadcastMessage(message)
	}
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
	writer.Write(serializedLog)
}

func listRooms(writer http.ResponseWriter, req *http.Request) {
	serializedRooms, err := json.Marshal(&rooms)
	if err != nil {
		fmt.Println("Marshalling the rooms has failed.")
		writer.WriteHeader(http.StatusInternalServerError)
	}
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
	found := false
	for _, user := range users {
		if user.UserName == userName {
			found = true
			roomsForUser = user.ChatRooms
		}
	}
	// if user not found, add the user
	if !found {
		users = append(users, &model.User{UserName: userName})
	}
	serializedRooms, err := json.Marshal(&roomsForUser)
	if err != nil {
		fmt.Println("Marshalling the rooms has failed.")
		writer.WriteHeader(http.StatusInternalServerError)
	}
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

	// If the room already exists, don't create a duplicate, return
	if _, err := getRoomForName(chatRoomName); err == nil {
		return
	}
	rooms = append(rooms, &model.ChatRoom{Users: nil, Name: chatRoomName, MaxUsers: maxUsers})
}

func joinRoom(writer http.ResponseWriter, req *http.Request) {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("Error reading the room from the request body")
		writer.WriteHeader(http.StatusInternalServerError)
	}
	var joinRequest *model.ChatRoomRequest
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
		// If the user isn't already in the room, add them to the room
		user := getUser(joinRequest.User.UserName)

		if room.GetUser(user.UserName) == nil {
			room.Users = append(room.Users, user)
		}
		user.Config = joinRequest.User.Config
		if user.GetRoom(room.Name) == nil {
			user.ChatRooms = append(user.ChatRooms, room)
		}

		room.Mux.Unlock()
		return
	}
}

func leaveRoom(writer http.ResponseWriter, req *http.Request) {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("Error reading the room from the request body")
		writer.WriteHeader(http.StatusInternalServerError)
	}
	var leaveRequest *model.ChatRoomRequest
	json.Unmarshal(bodyBytes, &leaveRequest)
	for _, room := range rooms {
		if room.Name != leaveRequest.RoomName {
			continue
		}

		// If the user is in the room, remove them from the room
		user := getUser(leaveRequest.User.UserName)

		room.Mux.Lock()
		// Update chat room
		if room.RemoveUser(user.UserName) != nil {
			room.Mux.Unlock()
			writer.WriteHeader(http.StatusForbidden)
			return
		}

		// Update user
		user.Config = leaveRequest.User.Config
		if user.RemoveRoom(room.Name) != nil {
			// Add the user back to the room, the operation failed
			room.Users = append(room.Users, user)
			room.Mux.Unlock()
			writer.WriteHeader(http.StatusForbidden)
			return
		}

		room.Mux.Unlock()
		return
	}
}

func updateUser(writer http.ResponseWriter, req *http.Request) {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("Error reading the message from the request body")
		writer.WriteHeader(http.StatusInternalServerError)
	}
	var user *model.User
	json.Unmarshal(bodyBytes, &user)
	updateUserConfig(user.UserName, user.Config)
}

// *************

func updateUserConfig(userName string, config *config.ClientConfig) {
	user := getUser(userName)
	if user != nil {
		user.Config = config
	} else {
		users = append(users, &model.User{UserName: userName, Config: config})
	}
}

func getUser(userName string) *model.User {
	// Find user
	var user *model.User
	for _, u := range users {
		if u.UserName == userName {
			user = u
		}
	}
	return user
}

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
	room, err := getRoomForName(message.ChatRoomName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, user := range room.Users {
		go sendMessageToUser(client, message, user)
	}
}

func sendMessageToUser(client *http.Client, message *model.Message, user *model.User) {
	// Format the message for serialization
	messageBuffer := model.ConvertMessageToBuffer(message)
	client.Post(fmt.Sprintf("http://localhost:%d/message", user.Config.MessagePort), "application/json; charset=utf-8", messageBuffer)
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
	http.HandleFunc("/user/update", updateUser)
	http.HandleFunc("/chatrooms/list", listRooms)
	http.HandleFunc("/chatrooms/create", createRoom)
	http.HandleFunc("/chatrooms/join", joinRoom)
	http.HandleFunc("/chatrooms/forUser", listRoomsForUser)
	http.HandleFunc("/chatrooms/leave", leaveRoom)

	start()
}
