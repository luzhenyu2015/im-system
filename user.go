package main

import "net"

type user struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

func NewUser(conn net.Conn, server *Server) *user {
	userAddr := conn.RemoteAddr().String()

	user := &user{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	go user.ListenMessage()
	return user
}

func (this *user) Online() {
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()
	this.server.Broadcast(this, "已上线")
}

func (this *user) Offline() {
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()
	this.server.Broadcast(this, "已下线")
}

func (this *user) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

func (this *user) DoMessage(msg string) {
	if msg == "who" {
		this.server.mapLock.RLock()
		for _, user := range this.server.OnlineMap {
			onlineMessage := "[" + user.Addr + "]" + user.Name + ":" + "在线。。\n"
			this.SendMsg(onlineMessage)
		}
		this.server.mapLock.RUnlock()
	} else {
		this.server.Broadcast(this, msg)
	}
}

func (this *user) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}
