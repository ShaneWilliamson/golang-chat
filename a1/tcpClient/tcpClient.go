package tcpClient

import (
	"436bin/a1/config"
	"436bin/a1/model"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

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

func (client *Client) readUserName() (string, error) {
	reader := createBufioReader()
	fmt.Print("Please enter your desired user name: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return text, nil
}

// CreateUser constructs the user with an associated UserName
func (client *Client) CreateUser() error {
	text, err := client.readUserName()
	if err != nil {
		return err
	}
	text = strings.TrimSpace(text)
	client.User.UserName = text
	return nil
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
	// Format the body for serialization
	joinRequest := &model.JoinChatRequest{User: client.User, RoomName: roomName}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(&joinRequest)

	resp, err := client.HTTPClient.Post("http://localhost:8081/chatrooms/join", "application/json; charset=utf-8", b)
	defer resp.Body.Close()
	if err != nil {
		// Failed to join the chat room
		DisplayErrorDialogWithMessage(err.Error())
		return
	}
	CreateChatTab(client, roomName)
}

func (client *Client) subscribeToServer() {
	fmt.Println("Starting client message subscription...")
	http.HandleFunc("/message", client.receiveMessage)
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
