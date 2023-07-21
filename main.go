package main

import (
	"flag"
	"fmt"
	"net"
	netcat "netcat/lib"
	"os"
	"time"
)

func main() {
	args := os.Args[1:]
	flag.IntVar(&netcat.MaxConnections, "m", 10, "Maximum number of concurrent connections allowed")
	flag.IntVar(&netcat.MaxLines, "l", 10000, "the max line of the log file allowed")
	flag.Parse()
	if len(args) > 1 {
		fmt.Println("❌ [USAGE]: ./TCPChat $port")
	} else {
		port := "8989"
		if len(args) != 0 {
			port = args[0]
		}
		listener, err := net.Listen("tcp", ":"+port)
		if err != nil {
			fmt.Println("Failed to launch server: ", err)
		}
		fmt.Printf("🚀 Server listening on the port :%s\n", port)

		timestamp := time.Now()
		fileName := fmt.Sprintf("logs/chat-log-%s.log", timestamp)
		netcat.LogFile, err = os.Create(fileName)
		if err != nil {
			fmt.Println("Failed to create log file: ", err)
		}
		for {
			connection, _ := listener.Accept()
			go netcat.Chat(connection)
		}
	}
}
