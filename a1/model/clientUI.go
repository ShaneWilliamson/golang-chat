package model

import (
	"github.com/therecipe/qt/widgets"
)

// ClientChatTab is used to keep track of information pertaining to the visually rendered parts of the tab
type ClientChatTab struct {
	Log          string
	LogView      *widgets.QTextEdit
	LogScrollBar *widgets.QScrollBar
}

// ChatRoomSelectionTab is used to keep track of information pertaining to the visually rendered parts of the tab
type ChatRoomSelectionTab struct {
}

// ClientUI keeps track of all tabs the client has
type ClientUI struct {
	ChatRoomSelectionTab *ChatRoomSelectionTab
	ChatTabs             *[]ClientChatTab
}
