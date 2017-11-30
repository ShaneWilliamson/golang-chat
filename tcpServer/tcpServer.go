package tcpServer

import (
	"container/heap"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"encoding/json"

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

func receiveMessage(message *model.Message) {
	room, err := getRoomForName(message.ChatRoomName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if room.GetUser(message.UserName) != nil {
		go logMessage(message)
		go broadcastMessage(message)
	}
	go HandleChatRoomDestruction()
}

func getLog(conn *net.Conn, chatRoomName string) {
	room, err := getRoomForName(chatRoomName)
	enc := json.NewEncoder(*conn)
	if err != nil {
		fmt.Println(err.Error())
		var errMessage []*model.Message
		errMessage = append(errMessage, model.ConstructErrorMessage(chatRoomName))
		enc.Encode(errMessage)
		return
	}
	enc.Encode(room.Log)
}

func listRooms(conn *net.Conn) {
	server := GetServerInstance()
	enc := json.NewEncoder(*conn)
	err := enc.Encode(&server.Rooms)
	if err != nil {
		fmt.Println("Listing the rooms to the client has failed.")
	}
}

func listRoomsForUser(conn *net.Conn, userName string) {
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
	enc := json.NewEncoder(*conn)
	err := enc.Encode(&roomsForUser)
	if err != nil {
		fmt.Println("Marshalling the rooms has failed.")
	}
}

func createRoom(chatRoomName string) {
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

func joinRoom(joinRequest *model.ChatRoomRequest) {
	server := GetServerInstance()
	for _, room := range server.Rooms {
		if room.Name != joinRequest.RoomName {
			continue
		}
		// Check, lock, check
		if len(room.Users) >= room.MaxUsers {
			fmt.Println("Cannot let user join, max users reached")
			return
		}
		room.Mux.Lock()
		if len(room.Users) >= room.MaxUsers {
			room.Mux.Unlock()
			fmt.Println("Cannot let user join, max users reached")
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

func leaveRoom(leaveRequest *model.ChatRoomRequest) {
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
			fmt.Println("Cannot remove user from room they are not in")
			return
		}

		// Update user
		user.Config = leaveRequest.User.Config
		if user.RemoveRoom(room.Name) != nil {
			// Add the user back to the room, the operation failed
			room.Users = append(room.Users, user)
			room.Mux.Unlock()
			fmt.Println("Removing the user from the room failed")
			return
		}

		room.Mux.Unlock()
		return
	}
}

// ***********************************

func createOrUpdateUser(u *model.User) {
	user := getUser(u.UserName)
	if user != nil {
		user.Config = u.Config
	} else {
		server := GetServerInstance()
		user = &model.User{UserName: u.UserName, Config: u.Config}
		server.Users = append(server.Users, user)
	}
	if user.Conn == nil {
		go connectToUser(user)
	}
}

func connectToUser(user *model.User) {
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", user.Config.MessagePort))
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Client dialing failed. Exiting.")
		os.Exit(1)
	}
	user.Conn = &conn
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
	room, err := getRoomForName(message.ChatRoomName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, user := range room.Users {
		go sendMessageToUser(message, user)
	}
}

func sendMessageToUser(message *model.Message, user *model.User) {
	// Get the user's connection if it exists
	if user.Conn == nil {
		return
	}
	req := model.ConvertToGenericRequest("Message", "message", nil, message)
	enc := json.NewEncoder(*user.Conn)
	enc.Encode(req)
}

func start(link net.Listener) {
	fmt.Println("Starting server...")
	// Listen for new client tcp socket connections
	for {
		conn := acceptNewClientConnection(link)
		go handleClientConnection(conn)
	}
}

func acceptNewClientConnection(link net.Listener) *net.Conn {
	// Wait for the next call, and returns a generic connection
	c, err := link.Accept()
	if err != nil {
		// Creation of new connection failed
		fmt.Println("Failed to accept a new connection. Exiting.")
		os.Exit(3)
	}
	return &c
}

func handleClientConnection(conn *net.Conn) {
	dec := json.NewDecoder(*conn)
	var user *model.User
	for {
		var req *model.GenericRequest
		err := dec.Decode(&req)
		if err == io.EOF {
			fmt.Println("Client connection closed")
			if user != nil {
				getUser(user.UserName).Conn = nil
			}
			return
		}
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		result := handleRequest(req, conn)
		if result != nil {
			user = result
		}
	}
}

func handleRequest(req *model.GenericRequest, conn *net.Conn) *model.User {
	data := model.ConvertFromGenericRequest(req)

	// ** Endpoints **
	switch req.Endpoint {
	case "log":
		getLog(conn, req.Params["roomName"])
	case "userUpdate":
		user, ok := data.(model.User)
		if !ok {
			fmt.Println("Failed to load User from request data")
			return nil
		}
		createOrUpdateUser(&user)
		return &user
	case "chatroomsList":
		listRooms(conn)
	case "chatroomsListForUser":
		listRoomsForUser(conn, req.Params["userName"])
	case "message":
		message, ok := data.(model.Message)
		if !ok {
			fmt.Println("Failed to load Message from request data")
			return nil
		}
		receiveMessage(&message)
	case "chatroomsCreate":
		createRoom(req.Params["roomName"])
	case "chatroomsJoin":
		chatReq, ok := data.(model.ChatRoomRequest)
		if !ok {
			fmt.Println("Failed to parse ChatRoomRequest from data")
			return nil
		}
		joinRoom(&chatReq)
	case "chatroomsLeave":
		chatReq, ok := data.(model.ChatRoomRequest)
		if !ok {
			fmt.Println("Failed to parse ChatRoomRequest from data")
			return nil
		}
		leaveRoom(&chatReq)
	}
	return nil
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
	// Get the user's connection if it exists
	if user.Conn == nil {
		return
	}
	req := model.ConvertToGenericRequest("User", "update", nil, user)
	enc := json.NewEncoder(*user.Conn)
	enc.Encode(req)
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
	// http.HandleFunc("/message", receiveMessage)
	// http.HandleFunc("/log", getLog)
	// http.HandleFunc("/user/update", updateUser)
	// http.HandleFunc("/chatrooms/list", listRooms)
	// http.HandleFunc("/chatrooms/create", createRoom)
	// http.HandleFunc("/chatrooms/join", joinRoom)
	// http.HandleFunc("/chatrooms/forUser", listRoomsForUser)
	// http.HandleFunc("/chatrooms/leave", leaveRoom)

	link, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Println("Error attempting to listen on port 8081. Exiting.")
		os.Exit(1)
	}

	start(link)
}
