package internal

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	MAX_CLIENTS   = 10
	INFOS         = ""
	CLIENTS_COUNT = 0
)

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

	err = os.Remove("data/history.txt")
	if err != nil {
		fmt.Println("No Previous chat history !")
	}

	err = os.WriteFile("data/history.txt", []byte(""), 0666)
	if err != nil {
		fmt.Println("Error creating chat history ...", err)
	}
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
			// connection.Write([]byte(label(name)))
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
		if INFOS != "" {
			logHistory(INFOS)
			INFOS = ""
		}
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

		message := label(name) + chatData
		logHistory(message)

		group.Range(func(key, value any) bool {
			// messageBrodcast(connection, value.(net.Conn), message)
			if value.(net.Conn) != connection && chatData != "\n" {
				infos := "\n" + label(name) + chatData + label(fmt.Sprintf("%s", key))
				value.(net.Conn).Write([]byte(infos))
				// value.(net.Conn).Write([]byte("YayY"))
				// } else {
			}
			return true
		})
		connection.Write([]byte(label(name)))

		if strings.TrimSpace(string(chatData)) == "STOP" || strings.TrimSpace(string(chatData)) == "EXIT" {
			fmt.Print(name + "Exiting Chat ... !")
			connection.Write([]byte("You've successfully logout !"))
			closeConnection(connection, group, name)
			// group.Delete(name)
			// connection.Close()
		}
	}
}

func messageBroadcast(group *sync.Map, connect net.Conn) {

}

// func newMessage(msg string, connection net.Conn) message {

// 	address := connection.RemoteAddr().String()
// 	return message{
// 		text:    address + msg,
// 		address: address,
// 	}
// }

func greeting() []byte {
	file, err := os.ReadFile("assets/welcome.txt")
	if err != nil {
		fmt.Println("Argh Welcome Message Missing !!!")
	}
	return file
}

func historyServing() []byte {
	file, err := os.ReadFile("data/history.txt")
	if err != nil {
		fmt.Println("Sorry chat history unavailable !")
		return []byte("Sorry chat history unavailable !\n")
	}
	if file == nil {
		return []byte("there is no chat history... \n")
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
	message := name + " has left our chat... \n"
	err := os.WriteFile("/data/hitory.txt", []byte(message), 0666)
	if err != nil {
		fmt.Println("Error logging chat history ...", err)
	}
	group.Range(func(key, value any) bool {
		if key != name {
			value.(net.Conn).Write([]byte("\n" + name + " has left our chat... \n"))
			value.(net.Conn).Write([]byte(label(fmt.Sprint(key))))
		}
		return true
	})
}

func joinMessage(conn net.Conn, group *sync.Map, name string) {

	// "\r\033[K"
	message := name + " has joined our chat... "
	INFOS = message + "\n"
	// err := os.WriteFile("/data/hitory.txt", []byte(message), 0666)
	// if err != nil {
	// }
	group.Range(func(key, value any) bool {
		message = "\n" + message + "\n" + label(fmt.Sprint(key))
		if key == name {
		} else {
			value.(net.Conn).Write([]byte(message))
		}
		return true
	})
}

func logHistory(message string) {
	file, err := os.OpenFile("data/history.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error with chat history file !", err)
		log.Println(err)
	}
	defer file.Close()
	if _, err := file.WriteString(message); err != nil {
		fmt.Println("Error logging chat history !", err)
		// log.Fatal(err)
	}
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
