package internal

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
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
	name    string
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
	fmt.Println("Server Ready on: ", PORT)

	var connectMap = &sync.Map{}
	// var connMap map[string]net.Conn

	for {
		connection, err := listener.Accept()
		connectionAddr := listener.Addr()
		if err != nil {
			fmt.Println(err)
			continue
		} //else {
		connection.Write(greeting())
		// name, _ := Client(&connection)
		access, name := accessChat(connection, connectionAddr, connectMap)
		fmt.Println(name + " joined")
		if access {
			connectMap.Store(name, connection)
			// connMap[name] = connection
			// fmt.Println(connectMap)
			go connectionHandler(connection, name, connectMap)
		}
	}
	// }

}

func connectionHandler(connection net.Conn, name string, group *sync.Map) {

	// defer closeConnection(connection, group)

	for {

		timing := timer()
		infos := timing + name + ":" // + string(chatData)
		// fmt.Print(infos)
		connection.Write([]byte(infos))
		// if strings.TrimSpace(s string)
		chatData, err := bufio.NewReader(connection).ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// fmt.Println(name + "Has LEFT !!!")
				closeConnection(connection, group, name)
				group.Delete(name)
				return
			} else {
				fmt.Println(err)
				return
			}
		}

		// message := chatData

		group.Range(func(key, value any) bool {
			// messageBrodcast(connection, value.(net.Conn), message)
			if value.(net.Conn) != connection {
				// value.(net.Conn).Write([]byte(message))
				// timing := timer()
				// infos := timing + fmt.Sprintf("%s", key) + ":"
				infos := "\r\033[K" + chatData + label(fmt.Sprintf("%s", key)) //timing + name + ":"
				value.(net.Conn).Write([]byte(infos))
			}
			return true
		})

		if strings.TrimSpace(string(chatData)) == "STOP" || strings.TrimSpace(string(chatData)) == "EXIT" {
			fmt.Print(name + "Exiting Chat ... !")
			connection.Write([]byte("You've successfully logout"))
			group.Delete(name)
			connection.Close()
		}

	}

}

func messageBroadcast(group *sync.Map, connect net.Conn) {

}

func newMessage(msg string, connection net.Conn) message {

	address := connection.RemoteAddr().String()
	return message{
		text:    address + msg,
		address: address,
	}
}

func greeting() []byte {
	file, err := os.ReadFile("assets/welcome.txt")
	if err != nil {
		fmt.Println("Argh Welcome Message Missing !!!")
	}
	return file
}

func label(name string) string {
	timing := timer()
	infos := timing + "[" + name + "]:" // + string(chatData)
	// fmt.Print(infos)
	return infos
}

func timer() string {
	t := time.Now()
	formatted := fmt.Sprintf("[%d-%02d-%02d %02d:%02d:%02d]",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	return formatted
}

func closeConnection(conn net.Conn, group *sync.Map, name string) {
	conn.Close()
	group.Delete(name)
	group.Range(func(key, value any) bool {
		if key != name {
			value.(net.Conn).Write([]byte(name + " has left our chat ... "))
		}
		return true
	})
}

func welcomeMessage(conn net.Conn, group *sync.Map, name string) {
	// conn.Close()
	// group.Delete(name)
	group.Range(func(key, value any) bool {
		if key == name {
		} else {
			value.(net.Conn).Write([]byte("\r\033[K" + name + " has joined our chat ... "))
		}
		return true
	})
}

// fmt.Println(*group)
// if key == name {
// 	return true
// } else {
// 	Write([]byte(message))
// return true

// })
// for {
// 	select {
// 	case msg := <-messages:
// 		for _, connect := range clients {
// 			if msg.address == connect.RemoteAddr().String() {
// 				continue
// 			}
// 			fmt.Fprintln(connect, msg.text)
// 		}
// 	case msg := <-leaving:
// 		for _, connect := range clients {
// 			fmt.Fprintln(connect, msg.text)
// 		}
// 	}
// }
// }

/*func leaveMessage(connect net.Conn, key, value interface{}, name string) {

}*/
