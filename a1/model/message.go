package model

// Message is what is transferred between the client and server, contains the username and body of the message
type Message struct {
	Sender string // TODO: make username-sender usage consistent, refactor
	Body   string
}

// // SerializeMessage attempts to convert a message object to a serializable format
// func SerializeMessage(m Message) string {
// 	// I'm thinking JSON style. There's probably some standard somewhere... todo?
// 	text := "{username: " + m.Sender + ", body: " + m.Body + "}\n"
// 	fmt.Println(text)
// 	return text
// }
