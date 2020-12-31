package base

import (
	"github.com/InoFlexin/serverbase/auth"
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
	Complex     bool
}

type Message struct {
	Json   string
	Key    string
	Action int
}

type SocketMessage struct {
	Packet Message
	Sock   net.Conn
}

type SocketEvent interface {
	OnMessageReceive(message *Message, client net.Conn)
	OnConnect(message *Message, client net.Conn)
	OnClose(err error)
}

var serverKey string = ""

func WriteSockMessage(socketMessage *SocketMessage) {
	socketMessage.Packet.Key = serverKey

	Write(&socketMessage.Packet, socketMessage.Sock)
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

func PacketMarshal(message *Message) []byte {
	e, err := json.Marshal(message)

	if err != nil {
		log.Println(err)
		return make([]byte, 0)
	}

	return e
}

func PacketUnmarshal(data []byte) *Message {
	message := Message{}
	json.Unmarshal([]byte(data), &message)

	return &message
}

func _receiveAndHandle(buf []byte, connection net.Conn, boot *Boot, count int) {
	var data = buf[:count]
	message := PacketUnmarshal(data)

	if message.Key != "" {
		switch message.Action {
		case ON_MSG_RECEIVE:
			go boot.Callback.OnMessageReceive(message, connection)
			break
		case ON_CONNECT:
			if GetKeyOrNil(message.Key) == "" {
				AddNewKey(message.Key)
				go boot.Callback.OnConnect(message, connection)
			}
			break
		}
	}
}

func receive(connection net.Conn, boot Boot, serverKey string) {
	buf := make([]byte, boot.ReceiveSize) //1kb

	for {
		count, error := connection.Read(buf)

		if nil != error {
			if io.EOF == error {
				log.Printf("connection is closed from client; %v", connection.RemoteAddr().String())
				RemoveKeyIfExsist(PacketUnmarshal(buf).Key) // if client connection refused, remove session.
				go boot.Callback.OnClose(nil)
			}

			log.Printf("fail to receive data; err: %v", error)
			return
		}

		if count > 0 {
			_receiveAndHandle(buf, connection, &boot, count)
		}
	}
}

func GetServerKey() string {
	return serverKey
}

func SetupComplexServer(boot Boot, wg *sync.WaitGroup) {
	if boot.Complex {
		wg.Done()
	}
}

func ServerStart(boot Boot, wg *sync.WaitGroup) {
	listener, error := net.Listen(boot.Protocol, boot.Port)
	serverKey = auth.GenerateKey(20)
	log.Println(boot)
	log.Println(boot.ServerName + " get started port: " + boot.Port)

	if error != nil {
		log.Fatalf("Failed to bind address to "+boot.Port+" err: %v", error)
	}

	defer listener.Close()
	defer wg.Done()
	SetupComplexServer(boot, wg)

	for {
		conn, error := listener.Accept()

		if nil != error {
			log.Printf("Failed to accept; err: %v", error)
			continue
		}

		go receive(conn, boot, serverKey)
	}
}
