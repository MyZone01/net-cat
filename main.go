package main

import (
	"fmt"
	"net"
	netcat "netcat/lib"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	args := os.Args[1:]
	if len(args) > 2 || (len(args) == 2 && !strings.HasPrefix(args[1], "-m=")) {
		fmt.Println("❌ [USAGE]: ./TCPChat $port")
		return
	} else {
		if len(args) == 2 {
			_maxConnection, err := strconv.Atoi(args[1][3:])
			if err != nil {
				fmt.Println("❌ [USAGE]: ./TCPChat $port -m=flag")
				return
			}
			netcat.MaxConnections = _maxConnection
		}
		port := "8989"
		if len(args) != 0 {
			port = args[0]
		}
		listener, err := net.Listen("tcp", ":"+port)
		if err != nil {
			fmt.Println("❌ Failed to launch server: ", err)
			return
		}
		fmt.Printf("🚀 Server listening on the port :%s\n", port)

		timestamp := time.Now()
		fileName := fmt.Sprintf("logs/chat-log-%s.log", timestamp)
		netcat.LogFile, err = os.Create(fileName)
		if err != nil {
			fmt.Println("❌ Failed to create log file: ", err)
			return
		}
		for {
			connection, err := listener.Accept()
			if err != nil {
				fmt.Println("❌ Can not connect to the server: ", err)
				return
			}
			go netcat.Chat(connection)
		}
	}
}
