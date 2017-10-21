
# TCP Socket Based Chat Client + Server

A client and server implementation which relies on HTTP for communication.


## Server Architecture

### Server actions:

* Create routes for the clients to hit with HTTP requests
  * Routes for chat rooms
    * Join, List, Create, Leave
  * Routes for Chat
    * Send message
* Broadcast new messages to clients

##### Messages

* Send message log for chosen chat room upon connection by client
* Receive message from user, and store into log
* Broadcast new messages received out to all connections to said room

