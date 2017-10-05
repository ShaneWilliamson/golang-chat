package tcpClient

import (
	"436bin/a1/model"
	"bytes"
	"encoding/gob"
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

func convertMessageToBuffer(message *model.Message) *bytes.Buffer {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(&message)
	return b
}

func printMessage(m *model.Message) {
	fmt.Printf("%s: %s\n", m.Sender, m.Body)
}

func (client *Client) sendMessage(text string) {
	// Format the message for serialization
	m := constructMessage(client.UserName, text)
	messageBuffer := convertMessageToBuffer(m)
	fmt.Printf("Message to send: %s\n", m) // todo: remove this

	resp, err := client.HTTPClient.Post("http://localhost:8081/message", "application/json; charset=utf-8", messageBuffer)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Successfully sent message")
}

func (client *Client) receiveMessage() (*model.Message, error) {
	message := &model.Message{}
	dec := gob.NewDecoder(*client.Conn)
	err := dec.Decode(message)
	if err != nil {
		fmt.Println("Decoding response from server failed.")
		return nil, err
	}

	return message, nil
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

// Create makes a new tcp client and waits to send a message to the target server.
func Create() {
	fmt.Println("Creating client...")
	// Create the client
	client := &Client{
		HTTPClient: &http.Client{},
	}

	// We need to create the connection first
	c, err := net.Dial("tcp", "127.0.0.1:8081")
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Client dialing failed. Exiting.")
		os.Exit(1)
	}
	client.Conn = &c

	// And now we create the GUI
	CreateChatWindow(client)
}
