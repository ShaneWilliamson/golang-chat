package model

// Reply is an generic rpc reply
type Reply struct {
	Value string
}

// ChatRoomRequest is how join/leave requests are shaped
type ChatRoomRequest struct {
	RoomName string
	User     *User
}
