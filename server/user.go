package main

import (
	"log"
	"net"
	"strings"
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

func (this *User) SendMsg(msg string) {
	this.Conn.Write([]byte(msg))
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
	if msg == "who" {
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Address + "]" + user.Name + ": " + "在线\n"
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]
		this.server.mapLock.Lock()
		if _, ok := this.server.OnlineMap[newName]; !ok {
			delete(this.server.OnlineMap, this.Name)
			this.Name = newName
			this.server.OnlineMap[newName] = this
			this.SendMsg("您已经更新用户名:" + this.Name + "\n")
		} else {
			this.SendMsg("用户名已被占用" + "\n")
		}
		this.server.mapLock.Unlock()
	} else if len(msg) > 4 && msg[:3] == "to|" {
		remoteName := strings.Split(msg, "|")[1]
		msg := "[" + this.Address + "]" + this.Name + ": " + strings.Split(msg, "|")[2] + " (私聊) \n"
		this.server.mapLock.Lock()
		if remoteUser, ok := this.server.OnlineMap[remoteName]; ok {
			remoteUser.SendMsg(msg)
		} else {
			this.SendMsg("该用户不存在")
		}
		this.server.mapLock.Unlock()

	} else {
		this.server.Broadcast(this, msg)
	}

}
