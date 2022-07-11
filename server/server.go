package main

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type Server struct {
	Address   string
	Port      int
	Message   chan string
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
}

func BuildServer(Address string, port int) *Server {
	server := &Server{
		Address:   Address,
		Port:      port,
		Message:   make(chan string),
		OnlineMap: make(map[string]*User),
	}
	return server
}

func (this *Server) handle(conn net.Conn) {
	user := NewUser(conn)
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()

	this.Broadcast(user, "上线")
	select {}
}

func (this *Server) Broadcast(user *User, msg string) {
	sendMsg := fmt.Sprintf("[%s]%s:%s", user.Address, user.Name, msg)
	this.Message <- sendMsg
}

func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message
		this.mapLock.Lock()
		for _, v := range this.OnlineMap {
			v.C <- msg
		}
		this.mapLock.Unlock()
	}
}

func (this *Server) start() {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Address, this.Port))
	if err != nil {
		log.Println(err)
		return
	}
	go this.ListenMessager()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
		}
		go this.handle(conn)
	}

}
