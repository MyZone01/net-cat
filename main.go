package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func main() {

	if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
	} else {
		PORT := ":"
		if len(os.Args) == 1 {
			PORT += "8989"
		} else {
			PORT += os.Args[1]
		}

		listener, err := net.Listen("tcp", PORT)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Choose between 1024 & 65535")
			return
		}
		defer listener.Close()

		connection, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		for {
			chatData, err := bufio.NewReader(connection).ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return
			}
			if strings.TrimSpace(string(chatData)) == "STOP" {
				fmt.Print("Exiting Chat ... !")
				return
			}

			// if strings.TrimSpace(s string)

			fmt.Print("->", string(chatData))
			timing := time.Now()
			timeInfos := timing.Format(time.RFC3339) + "\n"
			connection.Write([]byte(timeInfos))

		}
	}
}
