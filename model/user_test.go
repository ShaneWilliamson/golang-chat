package model

import (
	"fmt"
	"testing"
)

func TestGetRoomFindsExpectedRoom(t *testing.T) {
	expectedRoom := &ChatRoom{Name: "expected"}
	otherRoom := &ChatRoom{Name: "other"}
	user := &User{ChatRooms: []*ChatRoom{
		expectedRoom,
		otherRoom,
	}}

	// Execute block
	actualRoom := user.GetRoom(expectedRoom.Name)

	// Assert block
	if actualRoom == nil {
		fmt.Println("GetRoom did not find the expected room")
		t.Fail()
	}
	if actualRoom != expectedRoom {
		fmt.Printf("GetRoom did not return the expected room, instead it returned room: %s\n", actualRoom.Name)
		t.Fail()
	}
}

func TestGetRoomReturnsNilWhenNoRoomFound(t *testing.T) {
	expectedRoom := &ChatRoom{Name: "expected"}
	otherRoom := &ChatRoom{Name: "other"}
	user := &User{ChatRooms: []*ChatRoom{
		otherRoom,
	}}

	// Execute block
	actualRoom := user.GetRoom(expectedRoom.Name)

	// Assert block
	if actualRoom != nil {
		fmt.Printf("GetRoom found unexpected room: %s\n", actualRoom.Name)
		t.Fail()
	}
}

func TestRemoveRoomRemovesTheDesiredRoom(t *testing.T) {
	expectedRoom := &ChatRoom{Name: "expected"}
	otherRoom := &ChatRoom{Name: "other"}
	user := &User{ChatRooms: []*ChatRoom{
		expectedRoom,
		otherRoom,
	}}

	// Execute block
	err := user.RemoveRoom(expectedRoom.Name)

	// Assert block
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	if user.GetRoom(expectedRoom.Name) != nil {
		fmt.Println("Failed to remove the room from the user's chat rooms")
		t.Fail()
	}
}

func TestRemoveRoomReturnsErrorWhenRoomNotFound(t *testing.T) {
	otherRoom := &ChatRoom{Name: "other"}
	user := &User{ChatRooms: []*ChatRoom{
		otherRoom,
	}}

	// Execute block
	err := user.RemoveRoom("ThisDoesn'tExist")

	// Assert block
	if err == nil {
		fmt.Println("RemoveRoom did not return an error when trying to remove a non-existant room")
		t.Fail()
	}
	if err.Error() != "Could not find chat room" {
		fmt.Println(err.Error())
		t.Fail()
	}
}
