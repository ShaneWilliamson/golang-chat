package model

import (
	"bytes"
	"encoding/json"
)

// Message is what is transferred between the client and server, contains the username and body of the message
type Message struct {
	Sender string // TODO: make username-sender usage consistent, refactor
	Body   string
}

// ReadableFormat returns a human-friendly representation of the message
func (message *Message) ReadableFormat() string {
	return message.Sender + ": " + message.Body + "\n"
}

// ConvertMessageToBuffer takes a message and then returns it as a byte buffer
func ConvertMessageToBuffer(message *Message) *bytes.Buffer {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(&message)
	return b
}
