package base

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"sync"
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
	Action int
}

type SocketEvent interface {
	OnMessageReceive(message *Message, client net.Conn)
	OnConnect(message *Message, client net.Conn)
	OnClose(message *Message)
}

func Write(message *Message, client net.Conn) {
	//Encoding Message to Json String
	e, err := json.Marshal(message)

	if err != nil {
		log.Println(err)
		return
	}

	client.Write(e)
}

func PacketUnmarshal(data []byte) *Message {
	message := Message{}
	json.Unmarshal([]byte(data), &message)

	return &message
}

/*
	Please running goroutine
*/
func Broadcast(message *Message) {
	keys, values := GetSessions()
	length := len(keys)

	for i := 0; i < length; i++ {
		conn := values[i]

		Write(message, conn)
	}
}

func Receive(connection net.Conn, boot Boot) {
	buf := make([]byte, boot.ReceiveSize) //1kb
	var data []byte

	for {
		count, error := connection.Read(buf)

		if nil != error {
			if io.EOF == error {
				log.Printf("connection is closed from client; %v", connection.RemoteAddr().String())
				RemoveSession(connection.RemoteAddr().String()) // if client connection refused, remove session.
			}

			log.Printf("fail to receive data; err: %v", error)
			return
		}

		if count > 0 {
			data = buf[:count]

			message := PacketUnmarshal(data)
			go boot.Callback.OnMessageReceive(message, connection)
		}
	}
}

func ServerStart(boot Boot, wg *sync.WaitGroup) {
	listener, error := net.Listen(boot.Protocol, boot.Port)
	log.Println(boot)
	log.Println(boot.ServerName + " get started port: " + boot.Port)

	if error != nil {
		log.Fatalf("Failed to bind address to "+boot.Port+" err: %v", error)
	}
	defer listener.Close()
	defer boot.Callback.OnClose(&Message{Json: "close", Action: ON_CLOSE})
	wg.Done() //functions end calling wait group done()

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
