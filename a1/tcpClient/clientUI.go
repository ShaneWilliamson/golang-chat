package tcpClient

import (
	"436bin/a1/model"
	"fmt"

	"os"

	"github.com/therecipe/qt/widgets"
)

var log string
var logView *widgets.QTextEdit
var logScrollBar *widgets.QScrollBar
var tabWidget *widgets.QTabWidget

func addMessageToLogView(message *model.Message) {
	if logView == nil {
		fmt.Println("Log view not created yet, will not add new message to log")
		return
	}
	newEntry := message.ReadableFormat()
	log += newEntry
	logView.SetText(log)
	logScrollBar.SetValue(logScrollBar.Maximum())
}

func assignUserName(client *Client, username string) {
	// We assign the username to the client
	client.UserName = username
}

func createChatTab(client *Client) {
	tab := widgets.NewQWidget(nil, 0)
	layout := widgets.NewQVBoxLayout()
	logLayout := widgets.NewQVBoxLayout()
	inputLayout := widgets.NewQHBoxLayout()

	layout.InsertLayout(0, logLayout, 0)
	layout.InsertLayout(1, inputLayout, 0)
	tab.SetLayout(layout)

	logView = widgets.NewQTextEdit2("", nil)
	logView.SetReadOnly(true)
	logScrollBar = widgets.NewQScrollBar(nil)
	logView.SetVerticalScrollBar(logScrollBar)
	logLayout.AddWidget(logView, 0, 0)

	messageInput := widgets.NewQLineEdit(nil)
	messageInput.SetPlaceholderText("Enter message")
	submitButton := widgets.NewQPushButton2("Send", nil)
	fmt.Println(submitButton.AutoDefault())
	submitButton.SetAutoDefault(true)
	submitButton.ConnectClicked(func(checked bool) {
		client.sendMessage(messageInput.Text())
		// Reset the text
		messageInput.SetText("")
	})
	inputLayout.InsertWidget(0, messageInput, 0, 0)
	inputLayout.InsertWidget(1, submitButton, 0, 0)

	serverLog, err := getServerLog(client)
	if err != nil {
		logView.SetText("Unable to retrieve server log.")
	}
	// Set up the chat log
	log = ""
	for _, message := range serverLog {
		addMessageToLogView(message)
	}

	messageInput.ConnectReturnPressed(submitButton.Click)

	tabWidget.AddTab(tab, "chat")
}

// CreateChatWindow creates a window which contains the log and the ability to send messages
func CreateChatWindow(client *Client) {
	// Create application
	app := widgets.NewQApplication(len(os.Args), os.Args)

	// Create new tab widget
	tabWidget = widgets.NewQTabWidget(nil)

	// Create main window
	window := widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle("Hello World Example")
	window.SetMinimumSize2(200, 200)

	//********************************
	// Create main layout
	layout := widgets.NewQVBoxLayout()

	// Create main widget and set the layout
	mainWidget := widgets.NewQWidget(nil, 0)
	mainWidget.SetLayout(layout)

	// Create a line edit and add it to the layout
	usernameInput := widgets.NewQLineEdit(nil)
	usernameInput.SetPlaceholderText("Enter your username")
	layout.AddWidget(usernameInput, 0, 0)

	// Create a button and add it to the layout
	usernameButton := widgets.NewQPushButton2("Submit", nil)
	layout.AddWidget(usernameButton, 0, 0)
	//********************************

	// Connect event for button
	usernameButton.ConnectClicked(func(checked bool) {
		assignUserName(client, usernameInput.Text())
		tabWidget.RemoveTab(0)
		createChatTab(client)
	})

	tabWidget.AddTab(mainWidget, "Config")

	// Set main widget as the central widget of the window
	window.SetCentralWidget(tabWidget)

	// Show the window
	window.Show()

	// Execute app
	app.Exec()

	// err := ui.Main(func() {
	// 	entry := ui.NewEntry()
	// 	button := ui.NewButton("Ok")
	// 	logView = ui.NewLabel("")
	// 	box := ui.NewVerticalBox()
	// 	label := ui.NewLabel("Enter your name:")

	// 	// Adjust for initial username
	// 	box.Append(label, false)
	// 	box.Append(entry, false)
	// 	box.Append(button, false)
	// 	box.Append(logView, false)
	// 	window := ui.NewWindow("Hello", 200, 100, false)
	// 	window.SetChild(box)

	// 	// Maybe we make a new box here for the log and append the text as it comes in

	// 	button.OnClicked(func(*ui.Button) {
	// 		assignUserName(client, entry.Text())
	// 		// Reset the text
	// 		entry.SetText("")

	// 		// Username has been entered, now let's change this to send messages
	// 		label.SetText("Enter a message:")

	// 		// Retrieve log
	// serverLog, err := getServerLog(client)
	// if err != nil {
	// 	logView.SetText("Unable to retrieve server log.")
	// }

	// 		// Spin off the goroutine to update the log accordingly
	// 		for _, message := range serverLog {
	// 			addMessageToLogView(message)
	// 		}

	// 		// Now we make this button send messages
	// 		button.OnClicked(func(*ui.Button) {
	// 			// Send the message to the server
	// client.sendMessage(entry.Text())
	// // Reset the text
	// entry.SetText("")
	// 		})
	// 	})
	// 	window.OnClosing(func(*ui.Window) bool {
	// 		ui.Quit()
	// 		return true
	// 	})

	// 	window.Show()
	// })
	// if err != nil {
	// 	panic(err)
	// }
}
