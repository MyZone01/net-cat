package internal

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

var MAX_CONNECTION = 10
var count = 0
var clients = make(map[string]net.Conn)
var leaving = make(chan message)
var messages = make(chan message)

type Client struct {
	name    string
	address net.Conn
}

type message struct {
	text    string
	address string
}

func Server() {
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

	for len(clients) < 2 {

		connection, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		} //else {
		connection.Write(greeting())
		// name, _ := Client(&connection)
		access, name := accessChat(connection)
		if access {
			go connectionHandler(connection, name)
		}
	}
	// }

}

func connectionHandler(connection net.Conn, name string) {
	// connection.Write([]byte("[ENTER YOUR NAME]: "))

	for {
		timing := timer()
		infos := timing + name + ":" // + string(chatData)
		// fmt.Print(infos)
		connection.Write([]byte(infos))
		// if strings.TrimSpace(s string)
		chatData, err := bufio.NewReader(connection).ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("Has LEFT !!!")
				return
			} else {
				fmt.Println(err)
				return
			}
		}
		if strings.TrimSpace(string(chatData)) == "STOP" {
			fmt.Print("Exiting Chat ... !")
			return
		}

	}

}

func newMessage(msg string, connection net.Conn) message {
	address := connection.RemoteAddr().String()
	return message{
		text:    address + msg,
		address: address,
	}
}

func accessChat(connect net.Conn) (bool, string) {
	naming := ""
	for naming == "" {
		connect.Write([]byte("[ENTER YOUR NAME]: "))
		nami, err := bufio.NewReader(connect).ReadString('\n')
		if err != nil {
			if err == io.EOF {
				leaveMessage(connect)
				//fmt.Println("Has LEFT !!!")
				return false, naming
			} else {
				fmt.Println(err)
				return false, naming
			}
		}
		nami = strings.TrimSuffix(nami, "\n")
		naming = nami
		if err != nil {
			if err == io.EOF {
				fmt.Println("Joining room cancel !!!")
				// return
			} else {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		// goodPseudo:
		// connect.Write([]byte(greeting()))
		// connect.Write([]byte("[ENTER YOUR NAME]:"))
		// goto goodPseudo
	}
	name := "[" + naming + "]"
	return true, name
}

func greeting() []byte {
	file, err := os.ReadFile("assets/welcome.txt")
	if err != nil {
		fmt.Println("Argh Welcome Message Missing !!!")
	}
	return file
}

func timer() string {
	t := time.Now()
	formatted := fmt.Sprintf("[%d-%02d-%02d %02d:%02d:%02d]",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	return formatted
}

func messageBrodcast() {
	for {
		select {
		case msg := <-messages:
			for _, connect := range clients {
				if msg.address == connect.RemoteAddr().String() {
					continue
				}
				fmt.Fprintln(connect, msg.text)
			}
		case msg := <-leaving:
			for _, connect := range clients {
				fmt.Fprintln(connect, msg.text)
			}
		}
	}
}

func leaveMessage(connect net.Conn) {

}
