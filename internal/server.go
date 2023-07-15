package internal

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

var (
	MAX_CLIENTS = 10
	INFOS       = ""
	HISTORY     = "data/"
)

type CLIENTS struct {
	history string
	group   sync.Map
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

	err = os.RemoveAll("data")
	// err = os.Remove("data")
	if err != nil {
		fmt.Println(err)
	}

	if err = os.Mkdir("data", os.ModePerm); err != nil {
		log.Fatal(err)
	}
	fmt.Println("chatLogs directory successfully initialized !")

	fmt.Println("Server Ready on: ", PORT)

addGroup:

	// var connMap map[string]net.Conn
	// HISTORY = "data/history" +
	var groupie CLIENTS
	var connectMap = &sync.Map{}

	groupie.history = handlingLogs()

	for {
		connection, err := listener.Accept()
		connectionAddr := listener.Addr()
		if err != nil {
			fmt.Println(err)
			continue
		}
		connection.Write(greeting())

		// name, _ := Client(&connection)
		access, name := accessChat(connection, connectionAddr, connectMap, groupie.history)
		fmt.Println(name + " joined")
		if access {
			groupie.group = *connectMap
			connectMap.Store(name, connection)
			// connection.Write([]byte(label(name)))
			// connMap[name] = connection
			// fmt.Println(connectMap)
			go connectionHandler(connection, name, connectMap, groupie.history)
		}
		if checkGroup(connectMap) {
			goto addGroup
		}
	}

}

func checkGroup(group *sync.Map) bool {
	count := 0
	group.Range(func(key, value any) bool {
		count++
		return true
	})
	if count == MAX_CLIENTS {
		return true
	} else {
		return false
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

func historyServing(history string) []byte {
	file, err := os.ReadFile(history)
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

func closeConnection(conn net.Conn, group *sync.Map, name string, history string) {
	conn.Close()
	group.Delete(name)
	// message := name + " has left our chat... \n"
	// logHistory(message, history)
	// logHistory(message, history)
	// err := os.WriteFile(HISTORY, []byte(message), 0666)
	// if err != nil {
	// 	fmt.Println("Error logging chat history ...", err)
	// }
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
	// INFOS += message + "\n"
	// err := os.WriteFile("/data/hitory.txt", []byte(message), 0666)
	// if err != nil {
	// }
	// var container map[string]struct{}
	// container = map[string]struct{}{}
	group.Range(func(key, value interface{}) bool {
		msg := "\n" + message + "\n" + label(fmt.Sprint(key))
		if key == name {
		} else {
			connect := value.(net.Conn)
			connect.Write([]byte(msg))
		}
		return true
	})
}

func logHistory(message string, history string) {
	file, err := os.OpenFile(history, os.O_APPEND|os.O_WRONLY, 0644)
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

func handlingLogs() string {
	moment := time.Now()
	HISTORY += moment.String() + ".txt"
	err := os.WriteFile(HISTORY, []byte(""), 0666)
	if err != nil {
		fmt.Println("Error creating chat history ...", err)
	}
	return HISTORY
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
