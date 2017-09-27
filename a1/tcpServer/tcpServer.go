package tcpServer

import (
    "net"
    "errors"
    "fmt"
    "os"
    "bufio"
    "strings"
)


func acceptNewConn(link net.Listener) (net.Conn, error) {
    return link.Accept()
}


func readFromConnection(c net.Conn, link net.Listener) (string, error) {
    var retries = 10
    var message string
    var err error
    for i := 0; i < retries; i++ {
        message, err := bufio.NewReader(c).ReadString('\n')
        if err == nil {
            return message, nil
        }
        fmt.Println("Connection to client has been closed. Listening for new connection.")
        c, err = acceptNewConn(link)
        if err != nil {
            return nil, errors.New("Failed to listen for new connection.")
        }
    }
}


func run(c net.Conn, link net.Listener) {
    for {
		// will listen for message to process ending in newline (\n)
        message, err := readFromConnection(c)
        if err != nil {
            if err.Error() == "Failed to listen for new connection." {
                fmt.Println(err.Error())
                os.Exit(3)
            }
            fmt.Println("Failed to read, opening new connection.")
        }
        if err != nil {
            fmt.Println("Failed to read from connection, cannot recover. Exiting.")
            os.Exit(4)
        }
        // output message received
        fmt.Print("Message Received:", string(message))
        // sample process for string received
        newmessage := strings.ToUpper(message)
        // send new string back to client
        c.Write([]byte(newmessage + "\n"))
    }

}


// Create makes a new tcp server and listens for incoming requests
func Create() {
	// create the server
	fmt.Println("Creating Server...")
	link, err := net.Listen("tcp", ":8081")
    if err != nil {
		fmt.Println("Error attempting to listen on port 8081. Exiting.")
		os.Exit(1)
	}
	c, err := acceptNewConn(link)
    if err != nil {
		fmt.Println("Error connecting on port 8081")
		os.Exit(2)
	}
    run()
}

