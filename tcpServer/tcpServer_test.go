package tcpServer

import (
	"fmt"
	"testing"

	"github.com/ShaneWilliamson/golang-chat/config"
	"github.com/ShaneWilliamson/golang-chat/model"
)

func TestUpdateUserConfigCorrectlyUpdatesConfig(t *testing.T) {
	userName := "foo"
	user := &model.User{UserName: userName, Config: &config.ClientConfig{MessagePort: 1234}}
	updatedConfig := &config.ClientConfig{MessagePort: 9080}
	server := GetServerInstance()
	server.Users = []*model.User{
		user,
	}

	// Execute block
	updateUserConfig(userName, updatedConfig)

	// Assert block
	if user.Config.MessagePort != 9080 {
		fmt.Println("Failed to update the user's config")
		t.Fail()
	}
}

func TestUpdateUserConfigCorrectlyUpdatesConfigAfterCreatingUser(t *testing.T) {
	userName := "foo"
	updatedConfig := &config.ClientConfig{MessagePort: 9080}
	server := GetServerInstance()
	server.Users = []*model.User{}

	// Execute block
	updateUserConfig(userName, updatedConfig)

	// Assert block
	user := getUser(userName)
	if user.Config.MessagePort != 9080 {
		fmt.Println("Failed to update the user's config")
		t.Fail()
	}
}

func TestRemoveRoomFromServerCorrectlyRemoves(t *testing.T) {
	server := GetServerInstance()
	server.Rooms = []*model.ChatRoom{}
	targetRoom := &model.ChatRoom{Name: "TargetRoom"}
	server.Rooms = append(server.Rooms, &model.ChatRoom{Name: "Room1"})
	server.Rooms = append(server.Rooms, targetRoom)
	server.Rooms = append(server.Rooms, &model.ChatRoom{Name: "Room3"})

	// Execute block
	err := server.RemoveChatRoomFromRooms(targetRoom)

	// Assert block
	if err != nil {
		fmt.Printf(err.Error())
		t.Fail()
	}
	if len(server.Rooms) != 2 {
		fmt.Println("number of server rooms has not decreased")
		t.Fail()
	}
	// we know that there are 2 rooms, and they should be "Room1" and "Room3"
	if server.Rooms[0].Name != "Room1" || server.Rooms[1].Name != "Room3" {
		fmt.Printf("Rooms remaining in the server were not as expected. Index 0: %s, Index 1: %s\n", server.Rooms[0].Name, server.Rooms[1].Name)
		t.Fail()
	}
}

func TestRemoveRoomFromServerReturnsAnErrorWhenRoomDoesNotExist(t *testing.T) {
	server := GetServerInstance()
	server.Rooms = []*model.ChatRoom{}
	targetRoom := &model.ChatRoom{Name: "TargetRoom"}
	server.Rooms = append(server.Rooms, &model.ChatRoom{Name: "Room1"})
	server.Rooms = append(server.Rooms, &model.ChatRoom{Name: "Room3"})

	// Execute block
	err := server.RemoveChatRoomFromRooms(targetRoom)

	// Assert block
	if err == nil || err.Error() != "Could not find room by the name TargetRoom" {
		fmt.Printf("No error thrown when there should have been")
		t.Fail()
	}
	if len(server.Rooms) != 2 {
		fmt.Println("number of server rooms has changed when it should not have")
		t.Fail()
	}
	// we know that there are 2 rooms, and they should be "Room1" and "Room3"
	if server.Rooms[0].Name != "Room1" || server.Rooms[1].Name != "Room3" {
		fmt.Printf("Rooms remaining in the server were not as expected. Index 0: %s, Index 1: %s\n", server.Rooms[0].Name, server.Rooms[1].Name)
		t.Fail()
	}
}
