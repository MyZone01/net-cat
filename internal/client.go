package internal

//
import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
)

// checking if a name is provide and return a true that will allow access to the chat
func accessChat(connect net.Conn, connectAddr net.Addr, group *sync.Map, history string) (bool, string) {
	name := ""
	count := 0

	for name == "" || len(name) > 20 {
		if count == 3 {
			connect.Write([]byte("Please provide a non-empty name\n"))
		} else if count > 10 {
			connect.Write([]byte("You're a funny One !\n"))
		} else if count >= 50 {

		}
		connect.Write([]byte("[ENTER YOUR NAME]: "))
		count++
		naming, err := bufio.NewReader(connect).ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// leaveMessage(connect)
				Addr := connect.RemoteAddr().String()
				fmt.Println("Joining Cancelled ... from : " + Addr)
				return false, naming
			} else {
				fmt.Println(err)
				return false, naming
			}
		}

		name = strings.TrimSuffix(naming, "\n")
		if len(name) > 20 {
			connect.Write([]byte("Please Choose a smaller pseudo name"))
		}
	}
	// message := string(historyServing())
	joinMessage(connect, group, name)
	message := string(historyServing(history)) + label(name)
	connect.Write([]byte(message))
	// connect.Write([]byte(message))
	return true, name
}
func connectionHandler(connection net.Conn, name string, group *sync.Map, history string) {

	// defer closeConnection(connection, group)

	for {
		if INFOS != "" {
			logHistory(INFOS, history)
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
		logHistory(message, history)

		// var point CLIENTS
		// func (m *Map)
		// f := Load(key any) (value any, ok bool)
		// code := group
		group.Range(func(key, value any) bool {
			// val, ok := group.Load(key)
			// if ok {
			// messageBrodcast(connection, value.(net.Conn), message)
			// point := value.(net.Conn)
			if value.(net.Conn) != connection && chatData != "\n" {
				infos := "\n" + label(name) + chatData + label(fmt.Sprintf("%s", key))
				value.(net.Conn).Write([]byte(infos))
				// value.(net.Conn).Write([]byte("YayY"))
				// } else {
				// }
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
