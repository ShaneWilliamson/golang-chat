package tcpClient

import (
	"436bin/a1/model"
	"fmt"

	"github.com/andlabs/ui"
)

var window ui.Window
var logView *ui.Label

func addMessageToLogView(message *model.Message) {
	if logView == nil {
		fmt.Println("Log view not created yet, will not add new message to log")
		return
	}
	newEntry := message.ReadableFormat()
	ui.QueueMain(func() {
		logView.SetText(logView.Text() + newEntry)
	})
}

func assignUserName(client *Client, username string) {
	// We assign the username to the client
	client.UserName = username
}

// CreateChatWindow creates a window which contains the log and the ability to send messages
func CreateChatWindow(client *Client) {
	err := ui.Main(func() {
		entry := ui.NewEntry()
		button := ui.NewButton("Ok")
		logView = ui.NewLabel("")
		box := ui.NewVerticalBox()
		label := ui.NewLabel("Enter your name:")

		// Adjust for initial username
		box.Append(label, false)
		box.Append(entry, false)
		box.Append(button, false)
		box.Append(logView, false)
		window := ui.NewWindow("Hello", 200, 100, false)
		window.SetChild(box)

		// Maybe we make a new box here for the log and append the text as it comes in

		button.OnClicked(func(*ui.Button) {
			assignUserName(client, entry.Text())
			// Reset the text
			entry.SetText("")

			// Username has been entered, now let's change this to send messages
			label.SetText("Enter a message:")

			// Retrieve log
			serverLog, err := getServerLog(client)
			if err != nil {
				logView.SetText("Unable to retrieve server log.")
			}

			// Spin off the goroutine to update the log accordingly
			for _, message := range serverLog {
				addMessageToLogView(message)
			}

			// Now we make this button send messages
			button.OnClicked(func(*ui.Button) {
				// Send the message to the server
				client.sendMessage(entry.Text())
				// Reset the text
				entry.SetText("")
			})
		})
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})

		window.Show()
	})
	if err != nil {
		panic(err)
	}
}
