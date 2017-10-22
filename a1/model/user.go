package model

// User contains the username and their chatroom information
type User struct {
	UserName  string
	ChatRooms []*ChatRoom
}
