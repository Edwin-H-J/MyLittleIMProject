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
}

func NewUser(conn net.Conn) *User {
	user := &User{
		Address: conn.RemoteAddr().String(),
		Name:    conn.RemoteAddr().String(),
		C:       make(chan string),
		Conn:    conn,
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
