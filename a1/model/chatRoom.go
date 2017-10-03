package model

type ChatRoom struct {
	Users    *[]User
	Name     string
	MaxUsers int
}
