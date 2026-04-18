package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	OnlineMap map[string]*user
	mapLock   sync.RWMutex

	Message chan string
}

func NewServer(ip string, port int) *Server {
	Server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*user),
		Message:   make(chan string),
	}
	return Server
}

func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message

		this.mapLock.RLock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.RUnlock()
	}
}

func (this *Server) Broadcast(user *user, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
}

func (this *Server) Handle(conn net.Conn) {
	// 处理连接
	//fmt.Println("连接建立成功。。")
	user := NewUser(conn, this)
	user.Online()

	isLive := make(chan bool)

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("net.read err:", err)
				return
			}

			msg := string(buf[:n-1])
			user.DoMessage(msg)

			isLive <- true
		}
	}()

	for {
		select {
		case <-isLive:
		case <-time.After(time.Second * 10):
			user.SendMsg("你被踢了")
			close(user.C)
			conn.Close()
			return
		}
	}
}

func (this *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.listen err:", err)
		return
	}
	defer listener.Close()

	go this.ListenMessage()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("net.accept err:", err)
			continue
		}

		go this.Handle(conn)
	}

}
