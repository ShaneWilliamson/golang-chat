package tcpClient

import (
	"436bin/a1/model"
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
	UserName   string // TODO Change to be a User
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
	client.UserName = text
	return nil
}

// CreateReader this is split out for testing purposes
func createBufioReader() *bufio.Reader {
	return bufio.NewReader(os.Stdin)
}

func constructMessage(userName string, body string) *model.Message {
	// Marshal the message, and prepare it for transit
	message := &model.Message{Sender: userName, Body: strings.TrimSpace(body)}
	return message
}

func printMessage(m *model.Message) {
	fmt.Printf("%s: %s\n", m.Sender, m.Body)
}

func (client *Client) sendMessage(text string) {
	// Format the message for serialization
	m := constructMessage(client.UserName, text)
	messageBuffer := model.ConvertMessageToBuffer(m)

	resp, err := client.HTTPClient.Post("http://localhost:8081/message", "application/json; charset=utf-8", messageBuffer)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Successfully sent message")
}

func (client *Client) receiveMessage(writer http.ResponseWriter, req *http.Request) {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("Error reading the message from the request body")
		writer.WriteHeader(http.StatusBadRequest)
	}
	var message *model.Message
	json.Unmarshal(bodyBytes, &message)
	addMessageToLogView(message)
}

func readMessageFromUser(client *Client) (string, error) {
	fmt.Printf("%s: ", client.UserName)
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return text, nil
}

func getServerLog(c *Client) ([]*model.Message, error) {
	res, err := c.HTTPClient.Get("http://localhost:8081/log")
	// Wait for the response to complete
	defer res.Body.Close()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(bodyBytes))
	var serverLog []*model.Message
	json.Unmarshal(bodyBytes, &serverLog)
	return serverLog, nil
}

func subscribeToServer(client *Client) {
	fmt.Println("Starting client message subscription...")
	http.HandleFunc("/message", client.receiveMessage)
	fmt.Println((http.ListenAndServe(":9081", nil).Error()))
}

// Create makes a new tcp client and waits to send a message to the target server.
func Create() {
	fmt.Println("Creating client...")
	// Create the client
	client := &Client{
		HTTPClient: &http.Client{},
	}

	go subscribeToServer(client)

	// And now we create the GUI
	CreateChatWindow(client)
}
