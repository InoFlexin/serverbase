package main //simple server example

import (
	"fmt"

	"github.com/InoFlexin/serverbase/base"
)

type MyMessage base.Message

func (m MyMessage) OnMessageReceive(message *base.Message) {
	fmt.Println("on message receive: "+message.Json+" action: %d", message.Action)
}

func (m MyMessage) OnConnect(message *base.Message) {
	fmt.Println("on connect: "+message.Json+" action: %d", message.Action)
}

func (m MyMessage) OnClose(message *base.Message) {
	fmt.Printf("on close: "+message.Json+" action: %d", message.Action)
}

func main() {
	ev := MyMessage{}
	boot := base.Boot{Protocol: "tcp", Port: ":5092", ServerName: "test_server", Callback: ev}

	base.ServerStart(boot)
}
