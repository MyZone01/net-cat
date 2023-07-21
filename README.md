## 🌐 NET-CAT

TCPChat is a simple TCP-based chat server implemented in Go. It allows multiple clients to connect and communicate with each other in a real-time chat environment. The server keeps track of connected users, broadcasts messages to all connected users, and logs the chat history to a file.

## DESCRIPTION
The TCPChat project is a simple but effective implementation of a TCP-based chat server in Go. It allows users to connect, exchange messages, and see real-time updates from other connected users. The server provides a straightforward command-line interface and handles multiple connections concurrently using goroutines.

The project includes the following key features:

- **Real-time Communication**: Users can send and receive messages in real-time, creating an interactive chat experience.

- **User Management**: The server keeps track of connected users, ensuring unique names and limiting the number of simultaneous connections to 10.

- **Chat History Logging**: The server logs all chat messages to a file, enabling future reference and analysis of the conversation.

The TCPChat project can serve as a foundation for building more advanced chat applications or as a learning resource for understanding network programming concepts using Go.

We hope you find this project useful and enjoy exploring the world of TCP-based chat applications with TCPChat!

## USAGE

1. Download or clone the repository to your local machine.

2. Navigate to the project directory.

3. Build the server executable using the following command:
   ```
   go build .
   ```

4. Run the server using the following command:
   ```
   ./netcat
   ```

5. The server will start listening for incoming connections on the default port 8989. If you want to use a different port, you can specify it as a command-line argument:
   ```
   ./netcat 8080
   ```

6. Clients can connect to the server using Netcat client. For example, using Telnet:
   ```
   nc localhost 8989
   ```

7. Once connected, clients can start sending and receiving messages. Type your name to join the chat room.

## AVAILABLE COMMANDS

1. To clear the screen:
   ```
   /clear
   ```

2. To list all connected users:
   ```
   /list
   ```

3. To send a private message to a specific user:
   ```
   @username Your message here
   ```

4. To exit the chat:
   ```
   /exit
   ```

## NOTE

- The chat server has a maximum limit of 10 users. If the server is full, new connections will be rejected.

- The server logs all chat messages to a log file named "chat.log" in the "log" directory.

- Emoji shortcuts in messages will be automatically replaced with their corresponding emoji characters.

- Please ensure that you have Telnet or Netcat installed on your machine to connect to the server from the command line.

##  AUTHOR
+   Cheikh Ahmadou Tidiane Diallo
+   Serigne Saliou Mbacké Mbaye
+   Pape Bilaly Sow
