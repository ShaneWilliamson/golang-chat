package tcpServer

import (
	"bytes"
	"container/heap"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"encoding/json"

	"github.com/ShaneWilliamson/golang-chat/config"
	"github.com/ShaneWilliamson/golang-chat/model"
)

const sleepDays int = 7
const maxUsers int = 10

// Server is a singleton of the server
type Server struct {
	RoomDestructionMux sync.Mutex
	IsDestroying       bool
	Rooms              []*model.ChatRoom
	Users              []*model.User
}

var serverinstance *Server
var serveronce sync.Once

// GetServerInstance returns a singleton instance of the server
func GetServerInstance() *Server {
	serveronce.Do(func() {
		serverinstance = &Server{IsDestroying: false}
	})
	return serverinstance
}

// ***********************************

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
	go HandleChatRoomDestruction()
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
	server := GetServerInstance()
	serializedRooms, err := json.Marshal(&server.Rooms)
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
	server := GetServerInstance()
	for _, user := range server.Users {
		if user.UserName == userName {
			found = true
			roomsForUser = user.ChatRooms
		}
	}
	// if user not found, add the user
	if !found {
		server.Users = append(server.Users, &model.User{UserName: userName})
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
	pq := model.GetPriorityQueueInstance()
	chatRoom := &model.ChatRoom{Users: nil, Name: chatRoomName, MaxUsers: maxUsers, LastUsed: time.Now().UTC()}
	heap.Push(pq, chatRoom)
	server := GetServerInstance()
	server.Rooms = append(server.Rooms, chatRoom)
	go HandleChatRoomDestruction()
}

func joinRoom(writer http.ResponseWriter, req *http.Request) {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("Error reading the room from the request body")
		writer.WriteHeader(http.StatusInternalServerError)
	}
	var joinRequest *model.ChatRoomRequest
	json.Unmarshal(bodyBytes, &joinRequest)
	server := GetServerInstance()
	for _, room := range server.Rooms {
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
	server := GetServerInstance()
	for _, room := range server.Rooms {
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

// ***********************************

func updateUserConfig(userName string, config *config.ClientConfig) {
	user := getUser(userName)
	if user != nil {
		user.Config = config
	} else {
		server := GetServerInstance()
		server.Users = append(server.Users, &model.User{UserName: userName, Config: config})
	}
}

func getUser(userName string) *model.User {
	// Find user
	var user *model.User
	server := GetServerInstance()
	for _, u := range server.Users {
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
	// Update the room's LastUsed date, and it's order in the priority queue of destruction
	room.LastUsed = time.Now().UTC()
	pq := model.GetPriorityQueueInstance()
	pq.Update(room, room.LastUsed)
}

func getRoomForName(chatRoomName string) (*model.ChatRoom, error) {
	server := GetServerInstance()
	for _, room := range server.Rooms {
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

// HandleChatRoomDestruction launches a routine to destroy the next available chat room, returns if no rooms exist
func HandleChatRoomDestruction() {
	server := GetServerInstance()
	// If the server is already running a destruction sequence, return
	// Check, lock, check
	if server.IsDestroying {
		return
	}
	server.RoomDestructionMux.Lock()
	if server.IsDestroying {
		server.RoomDestructionMux.Unlock()
		return
	}
	server.IsDestroying = true
	// Get the priority queue of rooms to be destroyed
	pq := model.GetPriorityQueueInstance()
	for {
		// Get the last used room, sleep until its destruction
		if pq.Len() <= 0 { // Make sure there exists rooms
			server.IsDestroying = false
			server.RoomDestructionMux.Unlock()
			return
		}
		room := heap.Pop(pq).(*model.ChatRoom)
		// To ensure that the current target room is still the target after sleeping, add it back to the queue
		heap.Push(pq, room)
		sleepUntilTime := getSleepUntilTime(room)
		sleepDuration := getSleepDuration(sleepUntilTime)
		sleepUntil(sleepDuration)
		// Ensure that the initial target room is still our target room (It hasn't been used since sleeping)
		if pq.Len() <= 0 {
			// Somehow the queue has been emptied since we started sleeping, exit
			server.IsDestroying = false
			server.RoomDestructionMux.Unlock()
			return
		}
		currentTargetRoom := heap.Pop(pq).(*model.ChatRoom)
		sleepUntilTime = getSleepUntilTime(room)
		// If the time to which we want to sleep is after the current time (It has not yet passed) then continue
		if room.Name != currentTargetRoom.Name || sleepUntilTime.After(time.Now().UTC()) {
			// The room has been used since, we must sleep to the next target date
			heap.Push(pq, currentTargetRoom)
			continue
		}
		// Otherwise, if the time to which we want to sleep is before the current time, destroy the room
		destroyInactiveRoom(room)
	}
}

func sleepUntil(duration time.Duration) {
	time.Sleep(duration)
}

func getSleepUntilTime(room *model.ChatRoom) time.Time {
	// Calculate the duration to which we sleep
	return room.LastUsed.Add(time.Hour * time.Duration(sleepDays*24))
}

func getSleepDuration(sleepUntil time.Time) time.Duration {
	// Now, from the target date, calculate the duration from now.
	return sleepUntil.Sub(time.Now().UTC())
}

func destroyInactiveRoom(room *model.ChatRoom) {
	fmt.Printf("Destroying room %s", room.Name)
	// For every user in the room, remove the chat room from their array of rooms
	for _, user := range room.Users {
		user.RemoveRoom(room.Name)
		updateClient(user)
	}
	server := GetServerInstance()
	server.RemoveChatRoomFromRooms(room)
}

func updateClient(user *model.User) {
	client := &http.Client{}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(&user)
	client.Post(fmt.Sprintf("http://localhost:%d/update", user.Config.MessagePort), "application/json; charset=utf-8", b)
}

// RemoveChatRoomFromRooms removes the desired room from the server's array of rooms
func (server *Server) RemoveChatRoomFromRooms(room *model.ChatRoom) error {
	for i, r := range server.Rooms {
		if r.Name == room.Name {
			server.Rooms[len(server.Rooms)-1], server.Rooms[i] = server.Rooms[i], server.Rooms[len(server.Rooms)-1]
			server.Rooms = server.Rooms[:len(server.Rooms)-1]
			return nil
		}
	}
	return fmt.Errorf("Could not find room by the name %s", room.Name)
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
