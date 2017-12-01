package model

// Reply is an generic rpc reply
type Reply struct {
}

// ChatRoomRequest is how join/leave requests are shaped
type ChatRoomRequest struct {
	RoomName string
	User     *User
}
