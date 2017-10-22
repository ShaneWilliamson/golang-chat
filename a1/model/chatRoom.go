package model

import "sync"

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
