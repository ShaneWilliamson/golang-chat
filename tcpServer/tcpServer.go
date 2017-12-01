package tcpServer

import (
	"container/heap"
	"errors"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sync"
	"time"

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

// ReceiveMessage is an endpoint for receiving client messages
func (server *Server) ReceiveMessage(message *model.Message, reply *model.Reply) error {
	room, err := getRoomForName(message.ChatRoomName)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	if room.GetUser(message.UserName) != nil {
		go logMessage(message)
		go broadcastMessage(message)
	}
	go HandleChatRoomDestruction()
	*reply = model.Reply{Value: "foo"}
	return nil
}

// GetLog retrieves a specific chatroom log
func (server *Server) GetLog(roomName string, log *[]*model.Message) error {
	room, err := getRoomForName(roomName)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	*log = room.Log
	if len(*log) == 0 {
		return errors.New("No messages in log")
	}
	return nil
}

// ListRooms returns a list of all chatrooms to the user
func (server *Server) ListRooms(userName string, rooms *[]*model.ChatRoom) error {
	if userName == "" {
		*rooms = server.Rooms
		if len(*rooms) == 0 {
			return errors.New("No chatrooms exist")
		}
		return nil
	}
	*rooms = server.listRoomsForUser(userName)
	if len(*rooms) == 0 {
		return errors.New("No chatrooms exist")
	}
	return nil
}

func (server *Server) listRoomsForUser(userName string) []*model.ChatRoom {
	var roomsForUser []*model.ChatRoom
	found := false
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
	return roomsForUser
}

// CreateRoom creates a new chatroom
func (server *Server) CreateRoom(chatRoomName string, reply *model.Reply) error {
	// If the room already exists, don't create a duplicate, return
	if _, err := getRoomForName(chatRoomName); err == nil {
		return errors.New("Room already exists")
	}
	pq := model.GetPriorityQueueInstance()
	chatRoom := &model.ChatRoom{Users: nil, Name: chatRoomName, MaxUsers: maxUsers, LastUsed: time.Now().UTC()}
	heap.Push(pq, chatRoom)
	server.Rooms = append(server.Rooms, chatRoom)
	go HandleChatRoomDestruction()
	*reply = model.Reply{Value: "foo"}
	return nil
}

// JoinRoom joins a chatroom
func (server *Server) JoinRoom(req *model.ChatRoomRequest, reply *model.Reply) error {
	roomName := req.RoomName
	ru := req.User
	for _, room := range server.Rooms {
		if room.Name != roomName {
			continue
		}
		// Check, lock, check
		if len(room.Users) >= room.MaxUsers {
			return errors.New("Room full, cannot join")
		}
		room.Mux.Lock()
		if len(room.Users) >= room.MaxUsers {
			room.Mux.Unlock()
			return errors.New("Room full, cannot join")
		}
		// If the user isn't already in the room, add them to the room
		user := getUser(ru.UserName)

		if room.GetUser(user.UserName) == nil {
			room.Users = append(room.Users, user)
		}
		user.Config = ru.Config
		if user.GetRoom(room.Name) == nil {
			user.ChatRooms = append(user.ChatRooms, room)
		}

		room.Mux.Unlock()
		*reply = model.Reply{Value: "foo"}
		return nil
	}
	return errors.New("Could not find room to join")
}

// LeaveRoom leaves a chatroom
func (server *Server) LeaveRoom(req *model.ChatRoomRequest, reply *model.Reply) error {
	roomName := req.RoomName
	ru := req.User
	for _, room := range server.Rooms {
		if room.Name != roomName {
			continue
		}

		// If the user is in the room, remove them from the room
		user := getUser(ru.UserName)

		room.Mux.Lock()
		// Update chat room
		if err := room.RemoveUser(user.UserName); err != nil {
			room.Mux.Unlock()
			return err
		}

		// Update user
		user.Config = ru.Config
		if user.RemoveRoom(room.Name) != nil {
			// Add the user back to the room, the operation failed
			room.Users = append(room.Users, user)
			room.Mux.Unlock()
			return errors.New("Could not remove room from user")
		}

		room.Mux.Unlock()
		*reply = model.Reply{Value: "foo"}
		return nil
	}
	return nil
}

// UpdateUser updates a user's info
func (server *Server) UpdateUser(user *model.User, reply *model.Reply) error {
	updateUserConfig(user.UserName, user.Config)
	*reply = model.Reply{Value: "foo"}
	return nil
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
	// Format the message for serialization
	client, err := jsonrpc.Dial("tcp", fmt.Sprintf("localhost:%d", user.Config.MessagePort))
	if err != nil {
		log.Fatal("Server error:", err)
	}
	// Format the message for serialization
	var reply *model.Reply
	err = client.Call("RPCServer.ReceiveMessage", &message, &reply)
	if err != nil {
		log.Fatal("Server error:", err)
	}
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
	server.removeChatRoomFromRooms(room)
}

func updateClient(user *model.User) {
	// Format the message for serialization
	client, err := jsonrpc.Dial("tcp", fmt.Sprintf("localhost:%d", user.Config.MessagePort))
	if err != nil {
		log.Fatal("Server error:", err)
	}
	var reply *model.Reply
	err = client.Call("RPCServer.ReceiveUserUpdate", &user, &reply)
	if err != nil {
		log.Fatal("Server error:", err)
	}
}

// removeChatRoomFromRooms removes the desired room from the server's array of rooms
func (server *Server) removeChatRoomFromRooms(room *model.ChatRoom) error {
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

	// register for rpc
	server := GetServerInstance()
	s := rpc.NewServer()
	s.Register(server)
	s.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	listener, e := net.Listen("tcp", ":8081")
	if e != nil {
		log.Fatal("listener error:", e)
	}
	fmt.Println("Starting server...")
	for {
		if conn, err := listener.Accept(); err != nil {
			log.Fatal("accept error: " + err.Error())
		} else {
			log.Println("new connection established")
			go s.ServeCodec(jsonrpc.NewServerCodec(conn))
		}
	}
}
