package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ShaneWilliamson/golang-chat/tcpClient"
	"github.com/ShaneWilliamson/golang-chat/tcpServer"
)

func main() {
	// create server or clinet, and then loop until selection is made.
	reader := bufio.NewReader(os.Stdin)
	var text string
	for {
		fmt.Print("Please make a selection:\n1: Create server\n2: Create client\n\n(1/2): ")
		rawText, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error while reading choice please try again.")
			continue
		}
		text = strings.TrimSpace(rawText)
		if text != "1" && text != "2" {
			fmt.Println("Invalid selection, please select 1 for a new server or 2 for a new client.")
			fmt.Println(text)
			continue
		}
		break
	}
	switch text {
	case "1":
		tcpServer.Create()
	case "2":
		tcpClient.Create()
	default:
		fmt.Println("Invalid choice, unable to recover. Exiting.")
		os.Exit(1)
	}
}
