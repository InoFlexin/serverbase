package main //simple server example

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/InoFlexin/serverbase/auth"
	"github.com/InoFlexin/serverbase/base"
	"github.com/InoFlexin/serverbase/client"
)

//Override Server Message Type...
type (
	MyMessage       base.Message
	MyClientMessage base.Message
)

func (m MyMessage) OnMessageReceive(message *base.Message, client net.Conn) {
	fmt.Println("receive: " + message.Json + " key: " + message.Key)

	packetMessage := base.Message{Json: "pong", Key: base.GetServerKey(), Action: base.ON_MSG_RECEIVE}
	base.Write(&packetMessage, client)
}

func (m MyMessage) OnConnect(message *base.Message, client net.Conn) {
	fmt.Println("on  Connect!" + message.Key)
}

func (m MyMessage) OnClose(err error) {
	log.Println(err)
}

func (m MyClientMessage) OnMessageReceive(message *base.Message, server net.Conn) {
	fmt.Println("client receive: " + message.Json + " key: " + message.Key)
}

func (m MyClientMessage) OnConnect(message *base.Message, server net.Conn) {
	fmt.Println("client on connect: "+message.Json+" action: %d", message.Action)
}

func (m MyClientMessage) OnClose(err error) {
	log.Println(err)
}

func main() {
	wg := sync.WaitGroup{} //synchronized goroutine

	auth.RegisterKey("server", auth.GenerateKey(20))
	auth.RegisterKey("client", auth.GenerateKey(20))

	ev := MyMessage{}
	boot := base.Boot{Protocol: "tcp", Port: ":5092", ServerName: "test_server", Callback: ev, ReceiveSize: 1024, Complex: true}

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
