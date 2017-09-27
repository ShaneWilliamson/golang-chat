package tcpClient

import(
    "fmt"
    "net"
    "bufio"
    "os"
)

// Create makes a new tcp client and waits to send a message to the target server.
func Create() {
    c, err := net.Dial("tcp", "127.0.0.1:8081")
    if err != nil {
        fmt.Println("Client dialing failed. Exiting.")
        os.Exit(1)
    }
    for {
        reader := bufio.NewReader(os.Stdin)
        fmt.Print("sending message: ")
        text, err := reader.ReadString('\n')
        if err != nil {
            fmt.Print("Error while reading in message. Exiting.")
            os.Exit(1)
        }
        fmt.Fprintf(c, text + "\n")
        message, err := bufio.NewReader(c).ReadString('\n')
        if err != nil {
            fmt.Print("Failed to read response from server. Exiting.")
            os.Exit(2)
        }
        fmt.Println("Response from server:")
        fmt.Println(message)
    }
}


