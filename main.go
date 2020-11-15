package main //simple server example

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/InoFlexin/serverbase/base"
	"github.com/InoFlexin/serverbase/client"
)

type MyMessage base.Message
type MyClientMessage base.Message

func (m MyMessage) OnMessageReceive(message *base.Message, client net.Conn) {
	fmt.Println("on message receive: "+message.Json+" action: %d", message.Action)

	client.Write([]byte("pong"))
}

func (m MyMessage) OnConnect(message *base.Message, client net.Conn) {
	fmt.Println("on connect: "+message.Json+" action: %d", message.Action)
}

func (m MyMessage) OnClose(message *base.Message) {
	fmt.Printf("on close: "+message.Json+" action: %d", message.Action)
}

func (m MyClientMessage) OnMessageReceive(message *base.Message, server net.Conn) {
	fmt.Println("client on message receive: "+message.Json+" action: %d", message.Action)
}

func (m MyClientMessage) OnConnect(message *base.Message, server net.Conn) {
	fmt.Println("client on connect: "+message.Json+" action: %d", message.Action)
}

func (m MyClientMessage) OnClose(message *base.Message) {
	fmt.Printf("client on close: "+message.Json+" action: %d", message.Action)
}

func main() {
	wg := sync.WaitGroup{} //synchronized goroutine

	ev := MyMessage{}
	boot := base.Boot{Protocol: "tcp", Port: ":5092", ServerName: "test_server", Callback: ev, ReceiveSize: 1024}

	evm := MyClientMessage{}
	clientBoot := client.ClientBoot{Protocol: "tcp", HostAddr: "localhost", HostPort: ":5092", Callback: evm, BufferSize: 1024}

	wg.Add(1)
	go base.ServerStart(boot, &wg) //Server open
	wg.Wait()

	wg.Add(1)
	go client.ConnectServer(&clientBoot, &wg) //Client request connect to server
	wg.Wait()

	/*
		Client test logics
	*/
	serverError := client.SendPing(time.Second * 2)

	if serverError != nil {
		log.Fatal(serverError)
	}
}
