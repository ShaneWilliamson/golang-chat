package tcpClient

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/ShaneWilliamson/golang-chat/model"

	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// Client contains our username, and helper methods
type Client struct {
	User       *model.User
	HTTPClient *http.Client
	Conn       *net.Conn
	Link       *net.Listener
}

// CreateReader this is split out for testing purposes
func createBufioReader() *bufio.Reader {
	return bufio.NewReader(os.Stdin)
}

func constructMessage(chatRoomName string, userName string, body string) *model.Message {
	// Marshal the message, and prepare it for transit
	message := &model.Message{ChatRoomName: chatRoomName, UserName: userName, Body: strings.TrimSpace(body)}
	return message
}

func printMessage(m *model.Message) {
	fmt.Printf("%s: %s\n", m.UserName, m.Body)
}

func (client *Client) sendMessage(chatRoomName string, text string) {
	// Format the message for serialization
	m := constructMessage(chatRoomName, client.User.UserName, text)
	req := model.ConvertToGenericRequest("Message", "message", nil, m)
	enc := json.NewEncoder(*client.Conn)
	err := enc.Encode(req)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (client *Client) receiveMessage(message *model.Message) {
	var chatTab *model.ClientChatTab
	for _, room := range model.GetUIInstance().ChatTabs {
		if room.Name == message.ChatRoomName {
			chatTab = room
			break
		}
	}
	addMessageToLogView(message, chatTab)
}

func readMessageFromUser(client *Client) (string, error) {
	fmt.Printf("%s: ", client.User.UserName)
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return text, nil
}

func (client *Client) getServerLog(roomName string) ([]*model.Message, error) {
	params := make(map[string]string)
	params["roomName"] = roomName
	req := model.ConvertToGenericRequest("GetRequest", "log", params, nil)
	enc := json.NewEncoder(*client.Conn)
	dec := json.NewDecoder(*client.Conn)
	// Send the request
	enc.Encode(&roomName)
	// Wait for the response to complete
	var serverLog []*model.Message
	dec.Decode(&serverLog)
	return serverLog, nil
}

func (client *Client) getChatRooms() ([]*model.ChatRoom, error) {
	req := model.ConvertToGenericRequest("", "chatroomsList", nil, nil)
	enc := json.NewEncoder(*client.Conn)
	dec := json.NewDecoder(*client.Conn)
	// Send the request
	enc.Encode(&req)
	// Wait for the response to complete
	var chatRooms []*model.ChatRoom
	dec.Decode(&chatRooms)
	return chatRooms, nil
}

func (client *Client) getChatRoomsForUser() ([]*model.ChatRoom, error) {
	params := make(map[string]string)
	params["userName"] = client.User.UserName
	req := model.ConvertToGenericRequest("", "chatroomsListForUser", params, nil)

	enc := json.NewEncoder(*client.Conn)
	dec := json.NewDecoder(*client.Conn)
	// Send the request
	err := enc.Encode(&req)
	if err != nil {
		log.Fatal("Failed to send get request for chat rooms for user")
	}
	// Wait for the response to complete
	var chatRooms []*model.ChatRoom
	err = dec.Decode(&chatRooms)
	if err != nil {
		log.Fatal("Could not get chat rooms for user")
	}
	return chatRooms, nil
}

func (client *Client) createChatRoom(roomName string) {
	// Format the body for serialization
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(&roomName)

	resp, err := client.HTTPClient.Post("http://localhost:8081/chatrooms/create", "application/json; charset=utf-8", b)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("Successfully created room: %s\n", roomName)
}

func (client *Client) joinChatRoom(roomName string) {
	for _, tab := range model.GetUIInstance().ChatTabs {
		if tab.Name == roomName {
			return
		}
	}
	// Format the body for serialization
	joinRequest := &model.ChatRoomRequest{User: client.User, RoomName: roomName}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(&joinRequest)

	resp, err := client.HTTPClient.Post("http://localhost:8081/chatrooms/join", "application/json; charset=utf-8", b)
	defer resp.Body.Close()
	if err != nil || resp.StatusCode != 200 {
		// Failed to join the chat room
		DisplayErrorDialogWithMessage("Cannot join, max users reached")
		return
	}
	CreateChatTab(client, roomName)
}

func (client *Client) leaveChatRoom(roomName string) {
	// Format the body for serialization
	joinRequest := &model.ChatRoomRequest{User: client.User, RoomName: roomName}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(&joinRequest)

	resp, err := client.HTTPClient.Post("http://localhost:8081/chatrooms/leave", "application/json; charset=utf-8", b)
	defer resp.Body.Close()
	if err != nil || resp.StatusCode != 200 {
		// Failed to join the chat room
		DisplayErrorDialogWithMessage("Failed to leave chat room, please try again")
		return
	}
	tab, err := model.GetUIInstance().GetTabByName(roomName)
	if err != nil {
		DisplayErrorDialogWithMessage(err.Error())
		return
	}
	if model.GetUIInstance().RemoveTab(roomName) != nil {
		DisplayErrorDialogWithMessage("Critical error, could not remove client chat tab.")
		os.Exit(0)
	}
	model.GetUIInstance().TabWidget.RemoveTab(model.GetUIInstance().TabWidget.IndexOf(tab.Tab))
}

// UpdateUser creates/updates the user on the server
func (client *Client) UpdateUser() error {
	// Format the message for serialization

	model.ConvertToGenericRequest("User", "userUpdate", nil, &client.User)

	enc := json.NewEncoder(*client.Conn) // to write
	err := enc.Encode(&client.User)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (client *Client) receiveUserUpdate(user *model.User) {
	client.User = user
	model.GetUIInstance().ClientQTInstance.ReloadUI()
}

func (client *Client) handleRequests() {
	conn := client.listenForServer()
	dec := json.NewDecoder(*conn)
	for {
		var req *model.GenericRequest
		dec.Decode(req)
		client.handleGenericRequest(req)
	}
}

func (client *Client) handleGenericRequest(req *model.GenericRequest) {
	switch req.Endpoint {
	case "message":
		message, ok := model.ConvertFromGenericRequest(req).(*model.Message)
		if !ok {
			fmt.Println("Failed to receive message from server")
			return
		}
		client.receiveMessage(message)
	case "update":
		user, ok := model.ConvertFromGenericRequest(req).(*model.User)
		if !ok {
			fmt.Println("Failed to receive user update from server")
			return
		}
		client.receiveUserUpdate(user)
	default:
		log.Fatal("Invalid client endpoint")
	}
}

func (client *Client) listenForServer() *net.Conn {
	link, err := net.Listen("tcp", fmt.Sprintf(":%d", client.User.Config.MessagePort))
	if err != nil {
		fmt.Println("Error attempting to listen on client port. Exiting.")
		os.Exit(1)
	}
	var conn net.Conn
	conn, err = link.Accept()
	if err != nil {
		fmt.Println("Failed to accept connection from server, exiting")
		os.Exit(2)
	}
	return &conn
}

// ConnectToServer establishes an initial connection with the server, gives user info, and prompts for the user to establish further connection
func (client *Client) ConnectToServer() {
	// var err error
	// client.Link, err = net.Listen("tcp", Client.User.Config.MessagePort)
	// if err != nil {
	// 	fmt.Println("Error attempting to listen as client. Exiting.")
	// 	os.Exit(1)
	// }
	// Dial the server
	c, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Client dialing failed. Exiting.")
		os.Exit(1)
	}
	client.Conn = &c
	client.UpdateUser()
	go client.handleRequests()
}

// Create makes a new tcp client and waits to send a message to the target server.
func Create() {
	fmt.Println("Creating client...")

	// Create the client
	client := &Client{
		HTTPClient: &http.Client{},
		User:       &model.User{},
	}

	// And now we create the GUI
	chatApp := CreateChatWindow(client)
	chatApp.Exec()
}
