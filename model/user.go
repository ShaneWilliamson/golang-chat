package model

import (
	"errors"

	"github.com/ShaneWilliamson/golang-chat/config"
)

// User contains the username and their chatroom information
type User struct {
	UserName  string
	ChatRooms []*ChatRoom
	Config    *config.ClientConfig
}

// GetRoom returns the room if found, nil otherwise
func (user *User) GetRoom(roomName string) *ChatRoom {
	for _, r := range user.ChatRooms {
		if roomName == r.Name {
			return r
		}
	}
	return nil
}

// RemoveRoom removes the given chat room, returns an error if not found and removed
func (user *User) RemoveRoom(roomName string) error {
	for i, r := range user.ChatRooms {
		if roomName == r.Name {
			// Quick swap + remove chat room from chat rooms array
			user.ChatRooms[len(user.ChatRooms)-1], user.ChatRooms[i] = user.ChatRooms[i], user.ChatRooms[len(user.ChatRooms)-1]
			user.ChatRooms = user.ChatRooms[:len(user.ChatRooms)-1]
			return nil
		}
	}
	return errors.New("Could not find chat room")
}
