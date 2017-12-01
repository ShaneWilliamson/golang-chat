package tcpClient

import (
	"log"
	"net/http"
	"net/rpc"
	"sync"

	"github.com/ShaneWilliamson/golang-chat/config"
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
	RPCClient  *rpc.Client
	RPCServer  *RPCServer
	Conn       *net.Conn
}

var clientinstance *Client
var clientonce sync.Once

// GetClientInstance returns a singleton instance of the server
func GetClientInstance() *Client {
	clientonce.Do(func() {
		clientinstance = &Client{
			HTTPClient: &http.Client{},
			User:       &model.User{},
			RPCServer:  &RPCServer{},
		}
	})
	return clientinstance
}

// RPCServer is a placeholder struct as a client we use for providing RPC
type RPCServer struct {
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
	m := constructMessage(chatRoomName, client.User.UserName, text)
	var reply *model.Reply
	err := client.RPCClient.Call("Server.ReceiveMessage", m, &reply)
	if err != nil {
		log.Fatal("Server error:", err)
	}
}

func (s *RPCServer) receiveMessage(message *model.Message, reply *model.Reply) error {
	var chatTab *model.ClientChatTab
	for _, room := range model.GetUIInstance().ChatTabs {
		if room.Name == message.ChatRoomName {
			chatTab = room
			break
		}
	}
	addMessageToLogView(message, chatTab)
	*reply = model.Reply{}
	return nil
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
	var serverLog []*model.Message
	err := client.RPCClient.Call("Server.GetLog", roomName, &serverLog)
	if err != nil {
		log.Fatal("Server error:", err)
	}
	return serverLog, nil
}

func (client *Client) getChatRooms() ([]*model.ChatRoom, error) {
	var chatRooms []*model.ChatRoom
	// We don't want a specific user's rooms, so we don't provide a username
	err := client.RPCClient.Call("Server.ListRooms", "", &chatRooms)
	if err != nil {
		log.Fatal("Server error:", err)
	}
	return chatRooms, nil
}

func (client *Client) getChatRoomsForUser() ([]*model.ChatRoom, error) {
	var chatRooms []*model.ChatRoom
	// We want a specific user's rooms, so we provide a username
	err := client.RPCClient.Call("Server.ListRooms", client.User.UserName, &chatRooms)
	if err != nil {
		log.Fatal("Server error:", err)
	}
	return chatRooms, nil
}

func (client *Client) createChatRoom(roomName string) {
	var reply *model.Reply
	// We want a specific user's rooms, so we provide a username
	err := client.RPCClient.Call("Server.CreateRoom", roomName, &reply)
	if err != nil {
		log.Fatal("Server error:", err)
	}
	fmt.Printf("Successfully created room: %s\n", roomName)
}

func (client *Client) joinChatRoom(roomName string) {
	for _, tab := range model.GetUIInstance().ChatTabs {
		if tab.Name == roomName {
			return
		}
	}
	var reply *model.Reply
	// We want a specific user's rooms, so we provide a username
	err := client.RPCClient.Call("Server.JoinRoom", &model.ChatRoomRequest{RoomName: roomName, User: client.User}, &reply)
	if err != nil {
		// Failed to join the chat room
		DisplayErrorDialogWithMessage(err.Error())
		return
	}
	CreateChatTab(client, roomName)
}

func (client *Client) leaveChatRoom(roomName string) {
	// Format the body for serialization
	var reply *model.Reply
	// We want a specific user's rooms, so we provide a username
	err := client.RPCClient.Call("Server.LeaveRoom", &model.ChatRoomRequest{RoomName: roomName, User: client.User}, &reply)
	if err != nil {
		// Failed to join the chat room
		DisplayErrorDialogWithMessage(err.Error())
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

// UpdateUser updates the server about the user config
func (client *Client) UpdateUser() error {
	var reply *model.Reply
	// We want a specific user's rooms, so we provide a username
	err := client.RPCClient.Call("Server.JoinRoom", &client.User, &reply)
	if err != nil {
		log.Fatal("Server error:", err)
	}
	return nil
}

// ReceiveUserUpdate gets updates from the server when user data is updated
func (s *RPCServer) ReceiveUserUpdate(user *model.User, reply *model.Reply) error {
	GetClientInstance().User = user
	model.GetUIInstance().ClientQTInstance.ReloadUI()
	return nil
}

func (client *Client) subscribeToServer() {
	fmt.Println("Starting server subscription...")
	rpc.Register(client.RPCServer)
	rpc.HandleHTTP()
	var err error
	client.RPCClient, err = rpc.DialHTTP("tcp", "localhost:8081")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	fmt.Println((http.ListenAndServe(fmt.Sprintf(":%d", config.GetInstance().ClientConfig.MessagePort), nil).Error()))
}

// Create makes a new tcp client and waits to send a message to the target server.
func Create() {
	fmt.Println("Creating client...")
	// Create the client
	client := GetClientInstance()

	// And now we create the GUI
	chatApp := CreateChatWindow(client)
	chatApp.Exec()
}
