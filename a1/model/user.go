package model

import "436bin/a1/config"

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
