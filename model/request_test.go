package model

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenericRequestSuccessfullyUnmarshalsToMessage(t *testing.T) {
	message := &Message{
		ChatRoomName: "fooroom",
		UserName:     "barname",
		Body:         "bodybaz",
	}
	b, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	req := &GenericRequest{Type: "Message", Data: b}

	// Execute block
	result := ConvertFromGenericRequest(req)

	// Assert block
	_, ok := result.(Message)
	if !ok {
		fmt.Println("type assertion as Message failed")
		t.Fail()
	}
}

func TestGenericRequestSuccessfullyMarshalsFromMessage(t *testing.T) {
	message := &Message{
		ChatRoomName: "fooroom",
		UserName:     "barname",
		Body:         "bodybaz",
	}

	// Execute block
	req := ConvertToGenericRequest("Message", "", nil, message)

	// Assert block
	result := ConvertFromGenericRequest(req)
	_, ok := result.(Message)
	if !ok {
		fmt.Println("type assertion as Message failed")
		t.Fail()
	}
}

func TestGenericRequestHandlesNil(t *testing.T) {
	// Execute block
	req := ConvertToGenericRequest("Message", "", nil, nil)

	// Assert block
	result := ConvertFromGenericRequest(req)
	_, ok := result.(Message)
	if !ok {
		fmt.Println("type assertion as Message failed")
		t.Fail()
	}
}
