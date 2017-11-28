# Sockets

## Client
When client starts app, create a connection with the server.
When the client wants to perform any communication with the server, functions delegate their action to be performed using the connection.
### Conversion Tasks
* Get server log
* Send message
* Receive broadcasted messages
* Get list of chatrooms
* Get list of chatrooms for the user
* Create a chatroom
* Join a chatroom
* Leave a chatroom
* Send server client info on user-specified port, accept connection from server
### New Tasks
* Open client connection

## Server
When the server wants to broadcast messages utilize connection to communicate with the users.
When the server wants to update a user, utilize connections to communicate with the users.
### Conversion Tasks
* Broadcast messages to clients
* Update clients when rooms are deleted
* Get server to open connection with the user on user-specified port
### New Tasks
* Accept client connections

## Overall
### Tasks
* Struct with type and interface (pass generic data over wire)