package netcat

import (
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
	LogFile        *os.File
	MaxConnections = 10
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

// broadcast sends the given message to all connected users except the one with the specified name.
// It also logs the message to the LogFile if the lineCount is less than MaxLines.
// Parameters:
//   - message: The message to be broadcasted to the users.
//   - name: The name of the user who sent the message, used to exclude them from receiving the message.
//
// Goals:
//   - Send the message to all connected users.
//   - Log the message to the LogFile if the lineCount is less than MaxLines.
func broadcast(message, name string, mustLog bool) {
	if mustLog {
		_, err := LogFile.WriteString(message)
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
	if mustLog {
		messages += message
	}
}

// clearScreen clears the terminal screen for the given connection.
// It uses the "clear" command to clear the screen.
// Parameters:
//   - conn: A net.Conn representing the connection to the user's terminal.
func clearScreen(conn net.Conn) {
	cmd := exec.Command("clear")
	cmd.Stdout = conn
	cmd.Run()
}

// listUsers sends a list of users currently in the chat to the specified connection.
// Parameters:
//   - conn: A net.Conn representing the connection to the user.
func listUsers(conn net.Conn) {
	userList := "Users in chat:\n"
	for username := range users {
		userList += "- " + username + "\n"
	}
	conn.Write([]byte(userList))
}

// replaceEmojis replaces emoji shortcuts in the given text with their corresponding emoji characters.
// Parameters:
//   - text: A string containing text with emoji shortcuts to be replaced.
//
// Returns:
//   - A string with emoji shortcuts replaced by their corresponding emoji characters.
func replaceEmojis(text string) string {
	for k, v := range emojiMap {
		text = strings.Replace(text, k, v, -1)
	}
	return text
}
