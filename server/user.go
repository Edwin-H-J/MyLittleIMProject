package main

import (
	"log"
	"net"
)

type User struct {
	Address string
	Name    string
	C       chan string
	Conn    net.Conn
	server  *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	user := &User{
		Address: conn.RemoteAddr().String(),
		Name:    conn.RemoteAddr().String(),
		C:       make(chan string),
		Conn:    conn,
		server:  server,
	}
	go user.listenSend()
	return user
}

func (this *User) listenSend() {
	for {
		select {
		case msg := <-this.C:
			_, err := this.Conn.Write([]byte(msg + "\n"))
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func (this *User) Online() {
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	this.server.Broadcast(this, "上线")
}

func (this *User) Offline() {
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	this.server.Broadcast(this, "下线")
}
func (this *User) DoMessage(msg string) {
	this.server.Broadcast(this, msg)
}
