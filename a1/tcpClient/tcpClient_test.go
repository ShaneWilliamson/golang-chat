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
	if message.Body != "Bar" || message.Sender != "Foo" {
		fmt.Printf("Construction of message failed. Username: %s, Body: %s\n", message.Sender, message.Body)
		t.Fail()
	}
}
