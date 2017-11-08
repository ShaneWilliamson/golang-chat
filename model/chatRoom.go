package model

import (
	"errors"
	"sync"
)

// ChatRoom defines the structure of rooms which hold users that can converse
type ChatRoom struct {
	Users    []*User `json:"-"`
	Name     string
	MaxUsers int
	Mux      sync.Mutex
	Log      []*Message
}

// GetUser returns the user if found, nil otherwise
func (room *ChatRoom) GetUser(userName string) *User {
	for _, u := range room.Users {
		if userName == u.UserName {
			return u
		}
	}
	return nil
}

// RemoveUser removes the given user, returns an error if not found and removed
func (room *ChatRoom) RemoveUser(userName string) error {
	for i, u := range room.Users {
		if userName == u.UserName {
			// Quick swap + remove user from users array
			room.Users[len(room.Users)-1], room.Users[i] = room.Users[i], room.Users[len(room.Users)-1]
			room.Users = room.Users[:len(room.Users)-1]
			return nil
		}
	}
	return errors.New("Could not find user")
}
