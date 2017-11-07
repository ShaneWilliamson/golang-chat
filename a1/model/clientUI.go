package model

import (
	"fmt"
	"sync"

	"github.com/therecipe/qt/widgets"
)

// ClientChatTab is used to keep track of information pertaining to the visually rendered parts of the tab
type ClientChatTab struct {
	Log          string
	LogView      *widgets.QTextEdit
	LogScrollBar *widgets.QScrollBar
	Name         string
	Tab          *widgets.QWidget
}

// ChatRoomManagementTab is used to keep track of information pertaining to the visually rendered parts of the tab
type ChatRoomManagementTab struct {
	Tab *widgets.QWidget
}

// ClientUI keeps track of all tabs the client has
type ClientUI struct {
	ChatRoomManagementTab *ChatRoomManagementTab
	ChatTabs              []*ClientChatTab
	TabWidget             *widgets.QTabWidget
}

var instance *ClientUI
var once sync.Once

// GetUIInstance returns a singleton instance of the program configuration
func GetUIInstance() *ClientUI {
	once.Do(func() {
		instance = &ClientUI{ChatRoomManagementTab: &ChatRoomManagementTab{}}
	})
	return instance
}

// GetTabByName retrieves the tab by the given name, if it exists
func (ui *ClientUI) GetTabByName(tabName string) (*ClientChatTab, error) {
	for _, tab := range ui.ChatTabs {
		if tab.Name == tabName {
			return tab, nil
		}
	}
	return nil, fmt.Errorf("Could not find tab named %s", tabName)
}

// RemoveTab removes the tab by the given name, if it exists
func (ui *ClientUI) RemoveTab(tabName string) error {
	for i, tab := range ui.ChatTabs {
		if tab.Name == tabName {
			// Quick swap + remove chat room from chat rooms array
			ui.ChatTabs[len(ui.ChatTabs)-1], ui.ChatTabs[i] = ui.ChatTabs[i], ui.ChatTabs[len(ui.ChatTabs)-1]
			ui.ChatTabs = ui.ChatTabs[:len(ui.ChatTabs)-1]
			return nil
		}
	}
	return fmt.Errorf("Could not find tab named %s", tabName)
}
