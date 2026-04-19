package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *client {
	client := &client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return nil
	}

	client.conn = conn
	return client
}

func (client *client) DealResponse() {
	io.Copy(os.Stdout, client.conn)
}

func (client *client) menu() bool {
	var flag int

	fmt.Println("1. 公聊模式")
	fmt.Println("2. 私聊模式")
	fmt.Println("3. 更新用户名")
	fmt.Println("0. 退出")
	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("请输入合法范围内数字")
		return false
	}
}

func (client *client) updateName() bool {
	fmt.Println("请输入新用户名")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("net.Write err:", err)
		return false
	}
	return true
}

func (client *client) Run() {
	for client.flag != 0 {
		if client.menu() != true {

		}
		switch client.flag {
		case 1:
			fmt.Println("公聊模式选择..")
			break
		case 2:
			fmt.Println("私聊模式选择..")
			break
		case 3:
			client.updateName()
			break
		}
	}
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

	go client.DealResponse()

	fmt.Println("连接服务器成功")

	client.Run()
}
