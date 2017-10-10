package tcpClient

import (
	"fmt"
	"testing"
)

func TestReceiveErrorMessage(t *testing.T) {

}

func TestReceiveMessage(t *testing.T) {
	// todo
}

func TestConstructMessage(t *testing.T) {
	message := constructMessage("Foo", "Bar")
	if message.Body != "Bar" || message.UserName != "Foo" {
		fmt.Printf("Construction of message failed. Username: %s, Body: %s\n", message.UserName, message.Body)
		t.Fail()
	}
}
