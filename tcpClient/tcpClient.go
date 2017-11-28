package tcpClient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

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
	Conn       *net.Conn
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
	messageBuffer := model.ConvertMessageToBuffer(m)

	resp, err := client.HTTPClient.Post("http://localhost:8081/message", "application/json; charset=utf-8", messageBuffer)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func (client *Client) receiveMessage(writer http.ResponseWriter, req *http.Request) {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("Error reading the message from the request body")
		writer.WriteHeader(http.StatusBadRequest)
	}
	var message *model.Message
	json.Unmarshal(bodyBytes, &message)
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

	// Format the message for serialization
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(&roomName)
	res, err := client.HTTPClient.Post("http://localhost:8081/log", "application/json; charset=utf-8", b)
	// Wait for the response to complete
	defer res.Body.Close()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	var serverLog []*model.Message
	json.Unmarshal(bodyBytes, &serverLog)
	return serverLog, nil
}

func (client *Client) getChatRooms() ([]*model.ChatRoom, error) {
	res, err := client.HTTPClient.Get("http://localhost:8081/chatrooms/list")
	// Wait for the response to complete
	defer res.Body.Close()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	var chatRooms []*model.ChatRoom
	json.Unmarshal(bodyBytes, &chatRooms)
	return chatRooms, nil
}

func (client *Client) getChatRoomsForUser() ([]*model.ChatRoom, error) {
	// Format the body for serialization
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(&client.User.UserName)

	res, err := client.HTTPClient.Post("http://localhost:8081/chatrooms/forUser", "application/json; charset=utf-8", b)
	// Wait for the response to complete
	defer res.Body.Close()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	var chatRooms []*model.ChatRoom
	json.Unmarshal(bodyBytes, &chatRooms)
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

// UpdateUser updates the server about the user config
func (client *Client) UpdateUser() error {
	// Format the body for serialization
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(&client.User)

	res, err := client.HTTPClient.Post("http://localhost:8081/user/update", "application/json; charset=utf-8", b)
	// Wait for the response to complete
	defer res.Body.Close()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (client *Client) receiveUserUpdate(writer http.ResponseWriter, req *http.Request) {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("Error reading the message from the request body")
		writer.WriteHeader(http.StatusBadRequest)
	}
	var user *model.User
	json.Unmarshal(bodyBytes, &user)
	client.User = user
	model.GetUIInstance().ClientQTInstance.ReloadUI()
}

func (client *Client) subscribeToServer() {
	fmt.Println("Starting client message subscription...")
	http.HandleFunc("/message", client.receiveMessage)
	http.HandleFunc("/update", client.receiveUserUpdate)
	fmt.Println((http.ListenAndServe(fmt.Sprintf(":%d", config.GetInstance().ClientConfig.MessagePort), nil).Error()))
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
