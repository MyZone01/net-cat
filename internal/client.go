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
func accessChat(connect net.Conn, connectAddr net.Addr, group *sync.Map) (bool, string) {
	name := ""
	count := 0

	for name == "" || len(name) > 20 {
		if count == 3 {
			connect.Write([]byte("Please provide a non-empty name\n"))
		} else if count > 10 {
			connect.Write([]byte("You're a funny One !\n"))
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
	message := string(historyServing()) + label(name)
	connect.Write([]byte(message))
	// connect.Write([]byte(message))
	return true, name
}

func handlingGroup() {

}
