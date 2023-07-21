package netcat

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

// Chat handles the communication with a new user connecting to the chat.
// It checks if the chat is full, if not, it welcomes the user and asks for their name.
// If the name is already in use, it informs the user and asks for a different name.
// Once the user is added to the chat, they can participate in the conversation.
//
// Parameters:
//
//	user (net.Conn): The network connection representing the user.
//
// Goals:
//  1. Handle user registration and name validation.
//  2. Manage user interactions in the chat, such as sending messages and commands.
//  3. Handle special commands like renaming, private messages, listing users, and clearing the screen.
//  4. Broadcast messages to all other users in the chat.
//  5. Handle user disconnection and remove them from the chat.
func Chat(user net.Conn) {
	if len(users) >= MaxConnections {
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

	if name != "" {
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
}
