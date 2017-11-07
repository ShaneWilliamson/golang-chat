package tcpClient

import (
	"fmt"
	"testing"
)

func TestConstructMessage(t *testing.T) {
	message := constructMessage("Room1", "Foo", "Bar")
	if message.ChatRoomName != "Room1" || message.Body != "Bar" || message.UserName != "Foo" {
		fmt.Printf("Construction of message failed. Username: %s, Body: %s\n", message.UserName, message.Body)
		t.Fail()
	}
}
