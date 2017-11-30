package model

import (
	"encoding/json"
	"fmt"
	"log"
)

// ChatRoomRequest defines the structure of requests for leaving/joining a chat room
type ChatRoomRequest struct {
	RoomName string
	User     *User
}

// GenericRequest allows us to send generic json data across the wire to later decode
type GenericRequest struct {
	Type     string `json:type`
	Endpoint string
	Params   map[string]string
	Data     json.RawMessage
}

// ConvertToGenericRequest takes in an interface and marshals it into a GenericRequest
func ConvertToGenericRequest(interfaceType, endpoint string, params map[string]string, v interface{}) *GenericRequest {
	request := &GenericRequest{
		Type:     interfaceType,
		Endpoint: endpoint,
		Params:   params,
	}
	var err error
	request.Data, err = json.Marshal(v)
	if err != nil {
		log.Fatal(fmt.Sprintf("could not marshall %s into GenericRequest", interfaceType))
	}
	return request
}

// ConvertFromGenericRequest takes a GenericRequest and converts it to an appropriate object
func ConvertFromGenericRequest(request *GenericRequest) interface{} {
	switch request.Type {
	case "Message":
		var e Message
		if err := json.Unmarshal([]byte(request.Data), &e); err != nil {
			log.Fatal(err)
		}
		return e
	case "ChatRoomRequest":
		var e ChatRoomRequest
		if err := json.Unmarshal([]byte(request.Data), &e); err != nil {
			log.Fatal(err)
		}
		return e
	case "User":
		var e User
		if err := json.Unmarshal([]byte(request.Data), &e); err != nil {
			log.Fatal(err)
		}
		return e
	default:
		return nil
	}
}
