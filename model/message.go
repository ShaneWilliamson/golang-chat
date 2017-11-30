package model

import (
	"bytes"
	"encoding/json"
)

// Message is what is transferred between the client and server, contains the username and body of the message
type Message struct {
	ChatRoomName string
	UserName     string
	Body         string
}

// ReadableFormat returns a human-friendly representation of the message
func (message *Message) ReadableFormat() string {
	return message.UserName + ": " + message.Body + "\n"
}

// ConvertMessageToBuffer takes a message and then returns it as a byte buffer
func ConvertMessageToBuffer(message *Message) *bytes.Buffer {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(&message)
	return b
}

// ConstructErrorMessage returns a newly made error message, for a particular room
func ConstructErrorMessage(roomName string) *Message {
	return &Message{ChatRoomName: roomName, UserName: "SYSTEM", Body: "REQUEST ERROR"}
}
