package tcpClient

import (
	"436bin/a1/model"

	"github.com/andlabs/ui"
)

var window ui.Window

func addMessageToLogView(log *ui.Label, message *model.Message) {
	newEntry := message.ReadableFormat()
	ui.QueueMain(func() {
		log.SetText(log.Text() + newEntry)
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
		logView := ui.NewLabel("")
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
			log, err := getServerLog(client)
			if err != nil {
				logView.SetText("Unable to retrieve server log.")
			}

			// Spin off the goroutine to update the log accordingly
			for _, message := range log {
				addMessageToLogView(logView, message)
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

		// Spin up goroutine to handle populating the log view with new messages broadcasted from the server
		// TODO ...

		window.Show()
	})
	if err != nil {
		panic(err)
	}
}
