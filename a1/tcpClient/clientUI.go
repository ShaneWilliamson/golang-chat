package tcpClient

import (
	"436bin/a1/model"
	"fmt"

	"github.com/andlabs/ui"
)

var window ui.Window

func addMessageToLog(log *ui.Label, message *model.Message) {
	newEntry := message.ReadableFormat()
	ui.QueueMain(func() {
		log.SetText(log.Text() + newEntry)
	})
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
			// We assign the username to the client
			client.UserName = entry.Text()

			// Username has been entered, now let's change this to send messages
			label.SetText("Enter a message:")

			// Retrieve log
			log, err := getServerLog(client)
			if err != nil {
				fmt.Println("Panic?")
				panic(err) // todo, clean up
			}
			// split this out?
			logText := ""
			for _, logMessage := range log {
				logText += logMessage.ReadableFormat()
			}

			// Send the message to the server here
			logView.SetText(logText)

			// Spin off the goroutine to update the log accordingly
			go func() {
				for {
					message, err := client.receiveMessage()
					if err != nil {
						continue
					}
					addMessageToLog(logView, message)
				}
			}()

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
