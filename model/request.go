package model

// ChatRoomRequest defines the structure of requests for leaving/joining a chat room
type ChatRoomRequest struct {
	RoomName string
	User     *User
}
