package tcpClient

import (
	"436bin/a1/config"
	"436bin/a1/model"

	"fmt"
	"os"
	"strconv"

	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func addMessageToLogView(message *model.Message, chatTab *model.ClientChatTab) {
	if chatTab.LogView == nil {
		fmt.Println("Log view not created yet, will not add new message to log")
		return
	}
	newEntry := message.ReadableFormat()
	chatTab.Log += newEntry
	chatTab.LogView.SetText(chatTab.Log)
	scrollChatToBottom(chatTab)
}

// DisplayErrorDialogWithMessage shows an error dialog with the specified message
func DisplayErrorDialogWithMessage(errorMessage string) {
	widgets.QMessageBox_Information(nil, "Error", errorMessage,
		widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
}

func assignUserName(client *Client, username string) {
	// We assign the username to the client
	client.User.UserName = username
}

func scrollChatToBottom(chatTab *model.ClientChatTab) {
	chatTab.LogScrollBar.SetValue(chatTab.LogScrollBar.MaximumHeight())
}

// CreateChatTab creates a new chat tab
func CreateChatTab(client *Client, roomName string) {
	tab := widgets.NewQWidget(nil, 0)
	layout := widgets.NewQVBoxLayout()
	leaveLayout := widgets.NewQHBoxLayout()
	logLayout := widgets.NewQVBoxLayout()
	inputLayout := widgets.NewQHBoxLayout()

	layout.InsertLayout(0, leaveLayout, 0)
	layout.InsertLayout(1, logLayout, 0)
	layout.InsertLayout(2, inputLayout, 0)
	tab.SetLayout(layout)

	leaveButton := widgets.NewQPushButton2("Leave", nil)
	leaveButton.ConnectClicked(func(checked bool) {
		client.leaveChatRoom(roomName)
	})
	leaveLayout.InsertWidget(0, leaveButton, 0, 0)

	logView := widgets.NewQTextEdit2("", nil)
	logView.SetReadOnly(true)
	logScrollBar := widgets.NewQScrollBar(nil)
	logView.SetVerticalScrollBar(logScrollBar)
	logLayout.AddWidget(logView, 0, 0)

	messageInput := widgets.NewQLineEdit(nil)
	messageInput.SetPlaceholderText("Enter message")
	submitButton := widgets.NewQPushButton2("Send", nil)
	submitButton.SetAutoDefault(true)
	submitButton.ConnectClicked(func(checked bool) {
		if messageInput.Text() == "" {
			return
		}
		client.sendMessage(roomName, messageInput.Text())
		// Reset the text
		messageInput.SetText("")
	})
	inputLayout.InsertWidget(0, messageInput, 0, 0)
	inputLayout.InsertWidget(1, submitButton, 0, 0)

	serverLog, err := client.getServerLog(roomName)
	if err != nil {
		logView.SetText("Unable to retrieve server log.")
	}
	// Set up the chat log
	chatTab := &model.ClientChatTab{Log: "", LogView: logView, LogScrollBar: logScrollBar, Name: roomName, Tab: tab}
	for _, message := range serverLog {
		addMessageToLogView(message, chatTab)
	}

	messageInput.ConnectReturnPressed(submitButton.Click)

	model.GetUIInstance().TabWidget.AddTab(tab, roomName)
	model.GetUIInstance().ChatTabs = append(model.GetUIInstance().ChatTabs, chatTab)

	scrollChatToBottom(chatTab)
}

func createChatRoomSelectionTab(client *Client) {
	layout := createChatRoomSelectionLayout(client)

	tab := widgets.NewQWidget(nil, 0)
	tab.SetLayout(layout)

	model.GetUIInstance().TabWidget.InsertTab(0, tab, "Chat Rooms")
	model.GetUIInstance().ChatRoomManagementTab.Tab = tab
}

func createChatRoomSelectionLayout(client *Client) *widgets.QVBoxLayout {
	layout := widgets.NewQVBoxLayout()
	listJoinLayout := getChatRoomOptionsLayout(client)
	listRefreshLayout := getChatRoomRefreshLayout(client)
	createLayout := getChatRoomCreateLayout(client)

	layout.InsertLayout(0, listJoinLayout, 0)
	layout.InsertLayout(1, listRefreshLayout, 0)
	layout.InsertLayout(2, createLayout, 0)
	return layout
}

func getChatRoomOptionsLayout(client *Client) *widgets.QHBoxLayout {
	layout := widgets.NewQHBoxLayout()
	chatRooms, err := client.getChatRooms()
	if err != nil {
		errLabel := widgets.NewQLabel2("An error occurred", nil, 0)
		layout.InsertWidget(0, errLabel, 0, 0)
		return layout
	}
	if chatRooms == nil || len(chatRooms) == 0 {
		errLabel := widgets.NewQLabel2("No chat rooms exist", nil, 0)
		layout.InsertWidget(0, errLabel, 0, 0)
		return layout
	}
	roomsComboBox := makeChatRoomsComboBox(chatRooms)

	joinButton := widgets.NewQPushButton2("Join", nil)
	joinButton.ConnectClicked(func(checked bool) {
		client.joinChatRoom(roomsComboBox.CurrentText())
	})

	layout.InsertWidget(0, roomsComboBox, 0, 0)
	layout.InsertWidget(1, joinButton, 0, 0)

	return layout
}

func makeChatRoomsComboBox(chatRooms []*model.ChatRoom) *widgets.QComboBox {
	roomsComboBox := widgets.NewQComboBox(nil)
	var rooms []string
	for _, room := range chatRooms {
		rooms = append(rooms, room.Name)
	}
	roomsComboBox.AddItems(rooms)
	return roomsComboBox
}

func getChatRoomRefreshLayout(client *Client) *widgets.QVBoxLayout {
	refreshButton := widgets.NewQPushButton2("Refresh", nil)
	refreshButton.ConnectClicked(func(checked bool) {
		model.GetUIInstance().TabWidget.RemoveTab(model.GetUIInstance().TabWidget.IndexOf(model.GetUIInstance().ChatRoomManagementTab.Tab))
		createChatRoomSelectionTab(client)
		model.GetUIInstance().TabWidget.SetCurrentWidget(model.GetUIInstance().ChatRoomManagementTab.Tab)
	})

	layout := widgets.NewQVBoxLayout()
	layout.InsertWidget(0, refreshButton, 0, 2)

	return layout
}

func getChatRoomCreateLayout(client *Client) *widgets.QHBoxLayout {
	layout := widgets.NewQHBoxLayout()
	createInput := widgets.NewQLineEdit(nil)
	createInput.SetPlaceholderText("Enter new chat room name")
	submitButton := widgets.NewQPushButton2("Create", nil)
	submitButton.ConnectClicked(func(checked bool) {
		if createInput.Text() == "" {
			return
		}
		client.createChatRoom(createInput.Text())
		// Reset the text
		createInput.SetText("")
	})
	layout.InsertWidget(0, createInput, 0, 0)
	layout.InsertWidget(1, submitButton, 0, 0)

	return layout
}

// CreateChatWindow creates a window which contains the log and the ability to send messages
func CreateChatWindow(client *Client) *widgets.QApplication {
	// Create application
	app := widgets.NewQApplication(len(os.Args), os.Args)

	// Create new tab widget
	model.GetUIInstance().TabWidget = widgets.NewQTabWidget(nil)

	// Create main window
	window := widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle("Hello World Example")
	window.SetMinimumSize2(200, 200)

	//********************************
	// Create main layout
	layout := widgets.NewQVBoxLayout()
	//********************************
	// Create a line edit and add it to the layout
	portInput := widgets.NewQLineEdit(nil)
	portInput.SetPlaceholderText("Port")
	portInput.SetValidator(gui.NewQIntValidator(portInput))
	layout.AddWidget(portInput, 0, 0)

	// Create a button and add it to the layout
	portButton := widgets.NewQPushButton2("Submit", nil)
	layout.AddWidget(portButton, 0, 0)
	//********************************
	// Create main widget and set the layout
	mainWidget := widgets.NewQWidget(nil, 0)
	mainWidget.SetLayout(layout)

	// Create a line edit and add it to the layout
	usernameInput := widgets.NewQLineEdit(nil)
	usernameInput.SetPlaceholderText("Enter your username")
	usernameInput.SetEnabled(false)
	layout.AddWidget(usernameInput, 0, 0)

	// Create a button and add it to the layout
	usernameButton := widgets.NewQPushButton2("Submit", nil)
	usernameButton.SetEnabled(false)
	layout.AddWidget(usernameButton, 0, 0)
	//********************************

	// Connect event for button
	portButton.ConnectClicked(func(checked bool) {
		clientConfig := config.GetInstance()

		port, err := strconv.ParseUint(portInput.Text(), 10, 16)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		clientConfig.ClientConfig.MessagePort = uint16(port)
		client.User.Config = clientConfig.ClientConfig

		go client.subscribeToServer()

		usernameInput.SetEnabled(true)
		usernameButton.SetEnabled(true)
	})

	// Connect event for button
	usernameButton.ConnectClicked(func(checked bool) {
		// Don't allow empty strings
		if usernameInput.Text() == "" {
			return
		}
		assignUserName(client, usernameInput.Text())
		model.GetUIInstance().TabWidget.RemoveTab(0)
		createChatRoomSelectionTab(client)
		// Update the server about user config info
		err := client.UpdateUser()
		if err != nil {
			widgets.QMessageBox_Information(nil, "Error", "Failed to update server about client config",
				widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
			return
		}
		// Get chat rooms for user
		rooms, err := client.getChatRoomsForUser()
		if err != nil {
			widgets.QMessageBox_Information(nil, "Error", "Could not join user's chat rooms",
				widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
			return
		}
		if rooms == nil {
			return
		}
		for _, room := range rooms {
			CreateChatTab(client, room.Name)
		}
	})

	model.GetUIInstance().TabWidget.AddTab(mainWidget, "Config")

	// Set main widget as the central widget of the window
	window.SetCentralWidget(model.GetUIInstance().TabWidget)

	// Show the window
	window.Show()

	// Execute app
	return app
}
