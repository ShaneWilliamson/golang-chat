package model

// JoinChatRequest defines the structure of the request for joining a chat room
type JoinChatRequest struct {
	RoomName string
	User     *User
}
