package client

import (
	"io"
	"log"
	"net"

	"github.com/InoFlexin/serverbase/base"
)

type ClientBoot struct {
	Protocol   string
	HostAddr   string
	HostPort   string
	Callback   base.SocketEvent
	BufferSize uint64
}

func Write(json string, conn net.Conn) {
	conn.Write([]byte(json))
}

func Handle(boot *ClientBoot, conn net.Conn) {
	buf := make([]byte, boot.BufferSize) //1kb

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
			boot.Callback.OnMessageReceive(&base.Message{Json: string(data), Action: base.ON_MSG_RECEIVE}, conn)
		}
	}
}

func ConnectServer(boot *ClientBoot) {
	conn, err := net.Dial(boot.Protocol, boot.HostAddr+boot.HostPort)
	boot.Callback.OnConnect(&base.Message{Json: "-connect", Action: base.ON_CONNECT}, conn)

	if err != nil {
		log.Fatalf("faild to connec to server err %v", err)
	}
	defer conn.Close()
	defer boot.Callback.OnClose(&base.Message{Json: "-close", Action: base.ON_CLOSE})

	go Handle(boot, conn)
}
