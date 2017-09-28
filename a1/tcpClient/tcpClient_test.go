package tcpClient

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
)

var mockReaderObj mockReader

// implements the target interface
type mockReader struct {
	mock.Mock
}

func mockReadUserName(client *Client) (string, error) {
	return "Foobarbaz", nil
}

func TestCreateUser(t *testing.T) {
	expectedUserName := "Foobarbaz"
	// mockReaderObj := new(mockReader)
	// mockReaderObj.On("ReadString", '\n').Return(expectedUserName, nil)

	//Mock out the createReader func
	client := Client{readUserName: mockReadUserName}

	err := client.CreateUser()
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	if client.UserName != expectedUserName {
		fmt.Printf("Failed to create the user. Actual name: %s, Expected name: %s\n", client.UserName, expectedUserName)
		t.Fail()
	}
}

func TestConstructMessage(t *testing.T) {
	message := constructMessage("Foo", "Bar")
	if message.Body != "Bar" || message.Sender != "Foo" {
		fmt.Printf("Construction of message failed. Username: %s, Body: %s\n", message.Sender, message.Body)
		t.Fail()
	}
}
