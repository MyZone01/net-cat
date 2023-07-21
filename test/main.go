package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	var conn net.Conn
	var err error
	// Retry until connection is established
	for {
		conn, err = net.Dial("tcp", "11.11.90.185:2525")
		if err == nil {
			break
		}
		fmt.Println("❌ Failed to connect, retrying in 10 seconds:", err)
		time.Sleep(10 * time.Second)
	}
	defer conn.Close()

	fmt.Println("✅ Connection established")
	
	fmt.Fprintln(conn, "BotBot") // Sending the bot name to the server.
	
	fmt.Println("😂 Sending messages...")
	
	for i := 0; i < 100000; i++ { // Send 10000 messages.
		_, err := fmt.Fprintf(conn, "👋 Hello from BotBot %d\n", i)
		if err != nil {
			fmt.Println("Failed to send message:", err)
			os.Exit(1)
		}
		time.Sleep(1 * time.Millisecond) // Wait for 1 second before sending the next message.
	}

	fmt.Println("DONE.")
	// for i := 45; i < 254; i++ {
	// 	ip := fmt.Sprintf("11.11.90.%d:8989", i)
	// 	fmt.Println("Try to connect to:", ip)
	// 	conn, err := net.Dial("tcp", ip)
	// 	if err != nil {
	// 		fmt.Println("Failed to connect to ", ip, " : ", err)
	// 		// os.Exit(1)
	// 		continue
	// 	}
	// 	fmt.Println("Connection established with the ip:", ip)
	// 	defer conn.Close()

	// 	fmt.Fprintln(conn, "Bot-", i) // Sending the bot name to the server.

	// 	for i := 0; i < 10000; i++ { // Send 1000 messages.
	// 		_, err := fmt.Fprintf(conn, "This is message number %d\n", i)
	// 		if err != nil {
	// 			fmt.Println("Failed to send message:", err)
	// 			os.Exit(1)
	// 		}
	// 		time.Sleep(1 * time.Millisecond) // Wait for 1 millisecond before sending the next message.
	// 	}
	// }
}
