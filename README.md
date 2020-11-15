# Server, Client Framework - NWServerBaseFramework
서버, 클라이언트 개발을 EventListener 기반으로 쉽고 빠르게 개발 할 수있도록 만들어진 golang socket framework 입니다.  

# How to use server framework?
```go
package main

import (
    "fmt"
    "sync"

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
	wg := sync.WaitGroup{} 

	ev := MyMessage{} // 서버 이벤트 선언
	boot := base.Boot{Protocol: "tcp",
                      Port: ":5092",
                      ServerName: "test_server",
                      Callback: ev,
                      ReceiveSize: 1024}
	// server boot option 설정

	wg.Add(1) // synchronized gorutine
	go base.ServerStart(boot, &wg)
	wg.Wait()
}
```

# How to use client framework?
```go
package main

import (
	"fmt"
	"sync"

	"github.com/InoFlexin/serverbase/base"
	"github.com/InoFlexin/serverbase/client"
)

type MyClientMessage base.Message //Client Message 타입 정의

// =================== Client Event Listeners ======================
func (m MyClientMessage) OnMessageReceive(message *base.Message, server net.Conn) {
	fmt.Println("client on message receive: "+message.Json+" action: %d", message.Action)
}

func (m MyClientMessage) OnConnect(message *base.Message, server net.Conn) {
	fmt.Println("client on connect: "+message.Json+" action: %d", message.Action)
}

func (m MyClientMessage) OnClose(message *base.Message) {
	fmt.Printf("client on close: "+message.Json+" action: %d", message.Action)
}
// ==================================================================

func main() {
    wg := sync.WaitGroup{}

    event := MyClientMessage{}
    clientBoot := client.ClientBoot{Protocol: "tcp",
                                    HostAddr: "localhost",
                                    HostPort: ":5092",
                                    Callback: event, 
                                    BufferSize: 1024}		
    wg.Add(1) // synchronized goroutine
    go client.ConnectServer(&clientBoot, &wg)
    wg.Wait()

    /*
        해당 프레임워크에서는 클라이언트의 테스트를 위한 SendPing 함수가 존재한다.
    */
    serverError := client.SendPing(time.Second * 2)
    
    if serverError != nil {
            log.Fatal(serverError)
    }
}
```

# 개발자
남대영 - wsnam0507@gmail.com
