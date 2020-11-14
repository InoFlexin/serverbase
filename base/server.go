package base

import (
	"io"
	"log"
	"net"
)

type Boot struct {
	Protocol   string //server protocol type (tcp/ip, udp... etc)
	Port       string
	ServerName string
	Callback   SocketEvent
}

type Message struct {
	Json   string
	Action int //socket actions.
}

type SocketEvent interface {
	OnMessageReceive(message *Message)
	OnConnect(message *Message)
	OnClose(message *Message)
}

func Receive(connection net.Conn, event *SocketEvent) {
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
		}
	}

	msg := string(data)
	log.Println(msg)
}

func ServerStart(boot Boot) {
	listener, error := net.Listen(boot.Protocol, boot.Port)
	boot.Callback.OnConnect(&Message{Json: "connect", Action: ON_CONNECT})
	log.Println(boot.ServerName + " get started port: " + boot.Port)

	if error != nil {
		log.Fatalf("Failed to bind address to "+boot.Port+" err: %v", error)
	}
	defer listener.Close()

	for {
		conn, error := listener.Accept()

		if nil != error {
			log.Printf("Failed to accept; err: %v", error)
			continue
		}

		go Receive(conn, &boot.Callback)
	}

	boot.Callback.OnClose(&Message{Json: "close", Action: ON_CLOSE})
}
