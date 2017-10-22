package model

import (
	"sync"

	"github.com/therecipe/qt/widgets"
)

// ClientChatTab is used to keep track of information pertaining to the visually rendered parts of the tab
type ClientChatTab struct {
	Log          string
	LogView      *widgets.QTextEdit
	LogScrollBar *widgets.QScrollBar
}

// ChatRoomManagementTab is used to keep track of information pertaining to the visually rendered parts of the tab
type ChatRoomManagementTab struct {
	Tab *widgets.QWidget
}

// ClientUI keeps track of all tabs the client has
type ClientUI struct {
	ChatRoomManagementTab *ChatRoomManagementTab
	ChatTabs              []*ClientChatTab
}

var instance *ClientUI
var once sync.Once

// GetUIInstance returns a singleton instance of the program configuration
func GetUIInstance() *ClientUI {
	once.Do(func() {
		instance = &ClientUI{ChatRoomManagementTab: &ChatRoomManagementTab{}, ChatTabs: nil}
	})
	return instance
}
