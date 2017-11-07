package model

import (
	"fmt"
	"testing"
)

func TestGetUserFindsExpectedUser(t *testing.T) {
	expectedUser := &User{UserName: "foo"}
	user2 := &User{UserName: "bar"}
	room := &ChatRoom{
		Users: []*User{
			expectedUser,
			user2,
		},
	}

	// Execute block
	actualUser := room.GetUser(expectedUser.UserName)

	// Assert block
	if actualUser == nil {
		fmt.Println("GetUser did not find the expected user")
		t.Fail()
	}
	if actualUser != expectedUser {
		fmt.Printf("GetUser did not return the expected user, instead it returned user: %s\n", actualUser.UserName)
		t.Fail()
	}
}

func TestGetUserReturnsNilWhenNoUserFound(t *testing.T) {
	room := &ChatRoom{
		Users: []*User{
			&User{UserName: "foo"},
			&User{UserName: "bar"},
			&User{UserName: "baz"},
		},
	}

	// Execute block
	actualUser := room.GetUser("ThisNameDoesNotExist")

	// Assert block
	if actualUser != nil {
		fmt.Printf("GetUser found unexpected user: %s\n", actualUser.UserName)
		t.Fail()
	}
}

func TestRemoveUserRemovesTheDesiredUser(t *testing.T) {
	expectedUser := &User{UserName: "foo"}
	user2 := &User{UserName: "bar"}
	room := &ChatRoom{
		Users: []*User{
			expectedUser,
			user2,
		},
	}

	// Execute block
	err := room.RemoveUser(expectedUser.UserName)

	// Assert block
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	if room.GetUser(expectedUser.UserName) != nil {
		fmt.Println("Failed to remove the user from the chat room")
		t.Fail()
	}
}

func TestRemoveUserReturnsErrorWhenUserNotFound(t *testing.T) {
	room := &ChatRoom{
		Users: []*User{
			&User{UserName: "foo"},
			&User{UserName: "bar"},
			&User{UserName: "baz"},
		},
	}

	// Execute block
	err := room.RemoveUser("ThisUserDoesNotExist")

	// Assert block
	if err == nil {
		fmt.Println("RemoveUser did not return an error when user doesn't exist")
		t.Fail()
	}
	if err.Error() != "Could not find user" {
		fmt.Println(err.Error())
		t.Fail()
	}
}
