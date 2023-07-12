package main

import (
	"fmt"
	"netcat/internal"
	"os"
)

func main() {

	if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
	} else {
		internal.Server()
		internal.Client()
	}
}
