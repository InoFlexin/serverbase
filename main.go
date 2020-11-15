package main //simple server example

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/InoFlexin/serverbase/base"
)

type MyMessage base.Message

func (m MyMessage) OnMessageReceive(message *base.Message, client net.Conn) {
	fmt.Println("on message receive: "+message.Json+" action: %d", message.Action)

	base.Broadcast(&base.Message{Json: "Sended from servers", Action: base.ON_MSG_RECEIVE})
}

func (m MyMessage) OnConnect(message *base.Message, client net.Conn) {
	fmt.Println("on connect: "+message.Json+" action: %d", message.Action)
}

func (m MyMessage) OnClose(message *base.Message) {
	fmt.Printf("on close: "+message.Json+" action: %d", message.Action)
}

func Ping(server net.Conn) {
	for {
		server.Write([]byte("-ping"))
		time.Sleep(time.Second * 1)
	}
}

func Read(conn net.Conn) {
	buf := make([]byte, 1024) //1kb
	fmt.Println("Read from servers...")

	for {
		count, error := conn.Read(buf)

		if nil != error {
			if io.EOF == error {
				log.Printf("connection is closed from server; %v", conn.RemoteAddr().String())
			}

			log.Printf("fail to receive data; err: %v", error)
			return
		}

		if count > 0 {
			data := buf[:count]
			fmt.Println(string(data))
		}
	}
}

func TestClientConnection() {
	conn, err := net.Dial("tcp", ":5092")

	if err != nil {
		log.Fatalf("faild to connec to server err %v", err)
	}
	defer conn.Close()

	go Read(conn)
	Ping(conn)
}

func main() {
	ev := MyMessage{}
	boot := base.Boot{Protocol: "tcp", Port: ":5092", ServerName: "test_server", Callback: ev}

	go base.ServerStart(boot)
	time.Sleep(time.Second * 2)
	TestClientConnection()
}
