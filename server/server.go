package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
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
				fmt.Println("Conn Read err", err)
				return
			}
			msg := string(buf[0 : n-1])
			user.DoMessage(msg)
			isLive <- true
		}
	}()

	for {
		select {
		case <-isLive:
		case <-time.After(time.Second * 60):
			{
				user.SendMsg("你被踢了\n")
				close(user.C)
				conn.Close()
				return
			}
		}
	}
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
