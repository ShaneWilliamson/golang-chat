
# TCP Socket Based Chat Client + Server

A client and server implementation which relies on TCP sockets for communication.


## Server Architecture

### Lifecycle of the server:

* Server is created
* Server initializes an empty slice of Users (max 10 concurrent per room)
* Server creates pool of available connections
    * spin up 10 goroutines which wait for new connections
    * spin up x goroutines which wait for new connections
* The server waits for connections to close and spins up new goroutines

### Server Actions:

#### Chat rooms

* List chat rooms
* Create chat room
* Join chat room

#### Messages

* Send message log for chosen chat room upon connection by client
* Receive message from user, and store into log
* Broadcast new messages received out to all connections to said room

