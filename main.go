package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	users    = map[string]net.Conn{}
	messages string
	mutex    sync.Mutex
	logFile  *os.File
)

func broadcast(message, name string) {
	for i := range users {
		if i != name {
			users[i].Write([]byte("\r\033[K" + message))
		}
		users[i].Write([]byte(fmt.Sprintf("[%s][%s]:", time.Now().Format("2006-01-02 15:04:05"), i)))
	}
	messages += message
}

func chat(user net.Conn) {
	if len(users) == 10 {
		user.Write([]byte("🔒 Sorry chat is full. Try again later\n"))
		user.Close()
		return
	}
	var name, message string
	welcome, err := os.ReadFile("assets/logo.txt")
	if err != nil {
		fmt.Println("❌ [ERROR]: Cannot open file\n" + err.Error())
	}
	user.Write(welcome)
	scanner := bufio.NewScanner(user)
	for scanner.Scan() {
		name = scanner.Text()
		name = strings.TrimSpace(name)
		_user, isNameAlreadyUse := users[name]
		if len(name) == 0 {
			user.Write([]byte("❌ Please enter correct name.\n[ENTER YOUR NAME]:"))
		} else if isNameAlreadyUse {
			user.Write([]byte("❌ The name is already use in the room.\n[ENTER YOUR NAME]:"))
			_user.Write([]byte("\r\033[K ❌ Someone try to use your logging name\n"))
			_user.Write([]byte(fmt.Sprintf("[%s][%s]:", time.Now().Format("2006-01-02 15:04:05"), name)))
		} else {
			users[name] = user
			break
		}
	}

	user.Write([]byte(messages))
	broadcast(fmt.Sprintf("🤝 %s has joined our chat ...\n", name), name)
	for {
		ok := scanner.Scan()
		if !ok {
			break
		}
		text := scanner.Text()
		if len(strings.Trim(text, " \n\t\r")) != 0 {
			if strings.HasPrefix(text, "@") {
				_text := strings.Split(text, " ")
				private := string(_text[0][1:])
				message = fmt.Sprintf("[%s][%s]:", time.Now().Format("2006-01-02 15:04:05"), name) + strings.Join(_text[1:], " ") + "\n"
				users[private].Write([]byte("\r\033[K" + message + fmt.Sprintf("[%s][%s]:", time.Now().Format("2006-01-02 15:04:05"), private)))
				users[name].Write([]byte(fmt.Sprintf("[%s][%s]:", time.Now().Format("2006-01-02 15:04:05"), name)))
			} else {
				message = fmt.Sprintf("[%s][%s]:", time.Now().Format("2006-01-02 15:04:05"), name) + text + "\n"
				mutex.Lock()
				broadcast(message, name)
				_, err := logFile.WriteString(message)
				mutex.Unlock()
				if err != nil {
					fmt.Println("❌ [ERROR]: Cannot write on file " + err.Error())
				}
			}
		}
	}
	broadcast(fmt.Sprintf("👋 %s has left our chat ...\n", name), name)
	delete(users, name)
}

func main() {
	args := os.Args[1:]
	if len(args) > 1 {
		fmt.Println("❌ [USAGE]: ./TCPChat $port")
	} else {
		port := "8989"
		if len(args) != 0 {
			port = args[0]
		}
		listener, err := net.Listen("tcp", ":"+port)
		if err != nil {
			fmt.Println("Failed to launch server: ", err)
		}
		fmt.Printf("🚀 Server listening on the port :%s\n", port)
		logFile, err = os.Create("log/chat.log")
		if err != nil {
			fmt.Println("Failed to create log file: ", err)
		}
		for {
			connection, _ := listener.Accept()
			go chat(connection)
		}
	}
}
