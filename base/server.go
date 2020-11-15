package base

import (
	"io"
	"log"
	"net"
)

type Boot struct {
	Protocol    string //server protocol type (tcp/ip, udp... etc)
	Port        string
	ServerName  string
	Callback    SocketEvent
	ReceiveSize uint64
}

type Message struct {
	Json   string
	Action int //socket actions.
}

type SocketEvent interface {
	OnMessageReceive(message *Message, client net.Conn)
	OnConnect(message *Message, client net.Conn)
	OnClose(message *Message)
}

/*
	Please running goroutine
*/
func Broadcast(message *Message) {
	keys, values := GetSessions()
	length := len(keys)

	for i := 0; i < length; i++ {
		conn := values[i]
		id := keys[i]

		conn.Write([]byte(message.Json))

		log.Printf("Successfully sended data id: " + id)
	}
}

func Receive(connection net.Conn, boot Boot) {
	buf := make([]byte, 1024) //1kb
	var data []byte

	for {
		count, error := connection.Read(buf)

		if nil != error {
			if io.EOF == error {
				log.Printf("connection is closed from client; %v", connection.RemoteAddr().String())
			}

			log.Printf("fail to receive data; err: %v", error)
			return
		}

		if count > 0 {
			data = buf[:count]

			//TODO: Response 하는부분 설계해야함
			go boot.Callback.OnMessageReceive(&Message{Json: string(data), Action: ON_MSG_RECEIVE}, connection)
		}
	}
}

func ServerStart(boot Boot) {
	listener, error := net.Listen(boot.Protocol, boot.Port)

	log.Println(boot)
	log.Println(boot.ServerName + " get started port: " + boot.Port)

	if error != nil {
		log.Fatalf("Failed to bind address to "+boot.Port+" err: %v", error)
	}
	defer listener.Close()
	defer boot.Callback.OnClose(&Message{Json: "close", Action: ON_CLOSE})

	for {
		conn, error := listener.Accept()

		if nil != error {
			log.Printf("Failed to accept; err: %v", error)
			continue
		}

		id, sessionError := AddSession(conn.LocalAddr().String(), conn)
		boot.Callback.OnConnect(&Message{Json: "connect", Action: ON_CONNECT}, conn)

		if sessionError == nil {
			go Receive(conn, boot)
			log.Printf("Successfully session added id: " + id)
		} else {
			log.Fatal(sessionError)
		}
	}
}
