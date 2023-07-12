package internal

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func Client() {
	// if len(os.Args) == 1 {
	// fmt.Print("Please provide host:port")
	// } else {
	connect, err := net.Dial("tcp", os.Args[1])
	if err != nil {
		fmt.Println(err)
		fmt.Println("Choose between 1024 & 65535")
	}

goodName:
	nameRequest := bufio.NewReader(os.Stdin)
	fmt.Print("[ENTER YOUR NAME]:")
	name, err := nameRequest.ReadString('\n')
	if len(name) == 0 {
		goto goodName
	}

	for {
		greeting()
		timing := timer()
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		fmt.Fprint(connect, text+"\n")

		message, _ := bufio.NewReader(connect).ReadString('\n')
		fmt.Print(timing, "[", name, "]", message)
	}
	// }
}
