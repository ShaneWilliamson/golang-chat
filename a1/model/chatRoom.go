package model

import "sync"

// ChatRoom defines the structure of rooms which hold users that can converse
type ChatRoom struct {
	Users    []*User
	Name     string
	MaxUsers int
	Mux      sync.Mutex
	Log      []*Message
}
