
# TCP Socket Chat Client + Server

A client and server implementation which relies on TCP Sockets for communication.

![alt text](https://storage.googleapis.com/shanew/GolangChat/Chat1.png "Config Screen")

![alt text](https://storage.googleapis.com/shanew/GolangChat/Chat2.png "Chat Room Management Screen")

![alt text](https://storage.googleapis.com/shanew/GolangChat/Chat3.png "Chat Room Screen")

## Dependencies
* https://github.com/therecipe/qt

## Usage
1. Clone the repo and build the project using qt's toolset
    * `qtdeploy -docker build linux` (This assumes you're using the docker installation for qt)
2. Launch the server
    * `./deploy/linux/golang-chat.sh`
    * Select the server option (1)
    * The server utilizes port `8081`, make sure this is free
3. Launch the client(s)
    * `./deploy/linux/golang-chat.sh`
    * Select the client option (2)
    * **Note** every client must have a unique, free port they listen on, select unique ports
