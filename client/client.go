package client

import (
	"errors"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"../auth"
	"../base"
)

type ClientBoot struct {
	Protocol   string
	HostAddr   string
	HostPort   string
	Callback   base.SocketEvent
	BufferSize uint64
}

var server net.Conn = nil
var clientKey string = ""

func CreateError(errorMessage string) error {
	return errors.New(errorMessage)
}

func Write(json string) {
	message := &base.Message{Json: json, Key: clientKey, Action: base.ON_MSG_RECEIVE}
	log.Println("write: " + clientKey)
	e := base.PacketMarshal(message)
	log.Println("write2: " + string(e))

	server.Write(e)
}

func SendPing(duration time.Duration) error {
	var serverError error = nil

	if server != nil {
		for {
			Write("ping")
			time.Sleep(duration)
		}
	} else {
		serverError = CreateError("Server not connected error")
	}

	return serverError
}

func Handle(boot *ClientBoot, conn net.Conn, wg *sync.WaitGroup) {
	buf := make([]byte, boot.BufferSize)
	defer wg.Done() //핸들 처리용 wait group은 밑에 logic이 끝나면 Done 시킨다.

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
			boot.Callback.OnMessageReceive(base.PacketUnmarshal(data), conn)
		}
	}
}

func GetClientKey() string {
	return clientKey
}

func ConnectServer(boot *ClientBoot, wg *sync.WaitGroup) {
	conn, err := net.Dial(boot.Protocol, boot.HostAddr+boot.HostPort)
	server = conn
	clientKey = auth.GenerateKey(20)
	base.Write(&base.Message{Json: "connect!", Key: clientKey, Action: base.ON_CONNECT}, conn)

	defer conn.Close()

	if err != nil {
		log.Fatalf("faild to connected to server err %v", err)
	}

	wg.Done()                    //main에서 처리한 waitGroup은 done() 시킨다.
	handleWg := sync.WaitGroup{} //고루틴 핸들러 처리용 wait group
	handleWg.Add(1)
	go Handle(boot, server, &handleWg)
	handleWg.Wait()
}
