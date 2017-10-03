package model

import (
    ""
)

type ChatRoom struct {
    Users *[]User,
    Name string,
    MaxUsers int,
}

type User struct {
    UserName string,
}

