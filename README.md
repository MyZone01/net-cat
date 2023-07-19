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
* From the same computer:
1. Type the command into a new terminal (client's) to start:
```
$ nc localhost <port>
```
2. When connection is received, a linux logo would appear and ask for client's name
3. Enter your name and start typing!
4. For more than 1 client, open new terminals (maximum 10 connections), then use the same command and port to join the chat.
5. Start chatting!

##	OBSERVE
To observe TCP packet exchange while the program is running, you can use network monitoring tools. Here are a few options:

1. Wireshark: Wireshark is a popular network protocol analyzer that allows you to capture and inspect network packets. It supports various protocols, including TCP. You can run Wireshark on your machine and select the network interface through which the TCP packets are being exchanged. Wireshark will capture and display the packets, allowing you to analyze their content and observe the TCP packet exchange in real-time.

2. tcpdump: tcpdump is a command-line packet analyzer for capturing network traffic. It allows you to capture packets on a specific network interface and save them to a file for later analysis. You can run tcpdump in a separate terminal window while your program is running, specifying the network interface and filtering for TCP traffic. For example, the following command captures TCP packets on the "eth0" interface and saves them to a file called "packets.pcap":

   ```
   tcpdump -i eth0 -w packets.pcap tcp
   ```

   You can then open the "packets.pcap" file using Wireshark or another packet analyzer to inspect the captured TCP packets.

3. Netcat: If you want to observe the TCP packet exchange specifically for the program you mentioned earlier, you can use Netcat (nc) to create a listener on a specific port and redirect the program's output to that port. For example, if your program is running on port 8989, you can use the following command:

   ```
   nc -l 8989
   ```

   This will start a Netcat listener on port 8989, and any TCP packets sent by the program will be displayed in the terminal.

These tools allow you to monitor the TCP packet exchange between your program and other network entities. They provide valuable insights into the network communication and can help you debug and analyze network-related issues.

##  AUTHOR
+   Cheikh Ahmadou Tidiane Diallo
+   Serigne Saliou Mbacké Mbaye
+   Pape Bilaly Sow
