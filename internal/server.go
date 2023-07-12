package internal

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

var MAX_CONNECTION = 10
var count = 0
var clients = make(map[string]net.Conn)
var leaving = make(chan message)
var messages = make(chan message)

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

	for len(clients) < 10 {

		connection, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		} //else {
		// connection.Write([]byte(greeting()))
		access, name := accessChat(connection)
		if access {
			go connectionHandler(connection, name)
		}
		// }

	}
}

func connectionHandler(connection net.Conn, name string) {
	clients[connection.RemoteAddr().String()] = connection

	messages <- newMessage("joined.", connection)

	input := bufio.NewScanner(connection)
	for input.Scan() {
		messages <- newMessage(": "+input.Text(), connection)
	}

	delete(clients, connection.RemoteAddr().String())

	leaving <- newMessage("has left.", connection)

	connection.Close()
	/*
		for {
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

			// if strings.TrimSpace(s string)
			}*/

	timing := timer()
	infos := timing + name + ":" // + string(chatData)
	// fmt.Print(infos)
	connection.Write([]byte(infos))
}

func newMessage(msg string, connection net.Conn) message {
	address := connection.RemoteAddr().String()
	return message{
		text:    address + msg,
		address: address,
	}
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

func accessChat(connect net.Conn) (bool, string) {
	naming, err := bufio.NewReader(connect).ReadString('\n')
	if err != nil {
		if err == io.EOF {
			fmt.Println("Has LEFT !!!")
			os.Exit(0)
		} else {
			fmt.Println(err)
			os.Exit(1)
		}
	}
goodPseudo:
	connect.Write([]byte("[ENTER YOUR NAME]:"))
	if len(naming) == 0 {
		goto goodPseudo
	}
	return true, "[" + naming + "]"
}

func greeting() string {
	file, err := os.Open("assets/welcome.txt")
	if err != nil {
		fmt.Println("Argh Welcome Message Missing !!!")
	}
	defer file.Close()

	greeting := bufio.NewScanner(file)

	if err := greeting.Err(); err != nil {
		fmt.Println("You3're Welcome !!!")
	}
	// for greeting.Scan() {
	return greeting.Text()
	// }
}

func timer() string {
	t := time.Now()
	formatted := fmt.Sprintf("[%d-%02d-%02d %02d:%02d:%02d]",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	return formatted
}
