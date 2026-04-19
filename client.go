package main

import (
	"flag"
	"fmt"
	"net"
)

type client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

func NewClient(serverIp string, serverPort int) *client {
	client := &client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return nil
	}

	client.conn = conn
	return client
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "server ip")
	flag.IntVar(&serverPort, "port", 6666, "server port")
}

func main() {
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("连接服务器失败")
		return
	}
	fmt.Println("连接服务器成功")
	//select {}
}
