package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var (
	users          = map[string]net.Conn{}
	messages       string
	mutex          sync.Mutex
	logFile        *os.File
	maxConnections = 10
	maxLines       int
	lineCount      int
	emojiMap       = map[string]string{
		":)":  "\U0001F642", // slightly_smiling_face
		":(":  "\U0001F641", // slightly_frowning_face
		";)":  "\U0001F609", // wink
		"<3":  "\U0001F496", // red heart // yellow heart
		":D":  "\U0001F600", // grinning face
		":P":  "\U0001F61B", // face with tongue
		":O":  "\U0001F62E", // face with open mouth
		":'(": "\U0001F622", // crying face
		":/":  "\U0001F615", // confused face
		":*":  "\U0001F618", // face throwing a kiss
		";(":  "\U0001F622", // crying face
	}
)

func replaceEmojis(text string) string {
	for k, v := range emojiMap {
		text = strings.Replace(text, k, v, -1)
	}
	return text
}

func broadcast(message, name string) {
	// before writing to the log file, check if it exceeds the max line count
	if lineCount < maxLines {
		lineCount++
		_, err := logFile.WriteString(message)
		if err != nil {
			fmt.Println("❌ [ERROR]: Cannot write on file " + err.Error())
		}
	}
	for i := range users {
		if i != name {
			users[i].Write([]byte("\r\033[K" + message))
		}
		users[i].Write([]byte(fmt.Sprintf("[%s][%s]:", time.Now().Format("2006-01-02 15:04:05"), i)))
	}
	messages += message
}

func clearScreen(conn net.Conn) {
	// This clear command depends on the operating system
	// Here we're using a command that works in Unix systems (like Linux and MacOS)
	cmd := exec.Command("clear")
	cmd.Stdout = conn
	cmd.Run()
}

func listUsers(conn net.Conn) {
	userList := "Users in chat:\n"
	for username := range users {
		userList += "- " + username + "\n"
	}
	conn.Write([]byte(userList))
}

func chat(user net.Conn) {
	if len(users) >= maxConnections {
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
		text := replaceEmojis(scanner.Text())
		if len(strings.Trim(text, " \n\t\r")) != 0 {
			if strings.HasPrefix(text, ">rename ") {
				newName := strings.TrimSpace(text[8:])
				if newName == "" {
					user.Write([]byte("❌ Invalid name.\n"))
					user.Write([]byte(fmt.Sprintf("[%s][%s]:", time.Now().Format("2006-01-02 15:04:05"), name)))
				} else if _, ok := users[newName]; ok {
					user.Write([]byte("❌ Name is already in use.\n"))
					user.Write([]byte(fmt.Sprintf("[%s][%s]:", time.Now().Format("2006-01-02 15:04:05"), name)))
				} else {
					delete(users, name)
					prevName := name
					name = newName
					users[name] = user
					user.Write([]byte("✅ Name changed successfully.\n"))
					broadcast(fmt.Sprintf("🗣️  %s has change his name to %s ...\n", prevName, name), name)
				}
			} else if strings.HasPrefix(text, "@") {
				_text := strings.Split(text, " ")
				private := string(_text[0][1:])
				message = fmt.Sprintf("[%s][%s]:", time.Now().Format("2006-01-02 15:04:05"), name) + strings.Join(_text[1:], " ") + "\n"
				users[private].Write([]byte("\r\033[K" + message + fmt.Sprintf("[%s][%s]:", time.Now().Format("2006-01-02 15:04:05"), private)))
				users[name].Write([]byte(fmt.Sprintf("[%s][%s]:", time.Now().Format("2006-01-02 15:04:05"), name)))
			} else if text == ">list_user" {
				listUsers(user)
				user.Write([]byte(fmt.Sprintf("[%s][%s]:", time.Now().Format("2006-01-02 15:04:05"), name)))
			} else if text == ">clear" {
				clearScreen(user)
				user.Write([]byte(fmt.Sprintf("[%s][%s]:", time.Now().Format("2006-01-02 15:04:05"), name)))
			} else {
				message = fmt.Sprintf("[%s][%s]:", time.Now().Format("2006-01-02 15:04:05"), name) + text + "\n"
				mutex.Lock()
				broadcast(message, name)
				mutex.Unlock()
			}
		}
	}
	broadcast(fmt.Sprintf("👋 %s has left our chat ...\n", name), name)
	delete(users, name)
}

func main() {
	args := os.Args[1:]
	flag.IntVar(&maxConnections, "m", 10, "Maximum number of concurrent connections allowed")
	flag.IntVar(&maxLines, "l", 10000, "the max line of the log file allowed")
	flag.Parse()
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
