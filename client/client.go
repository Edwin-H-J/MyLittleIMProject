package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}
	client.conn = conn
	return client
}

func (this *Client) menu() bool {
	var flag = 999
	fmt.Println("1.群发模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")
	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		this.flag = flag
		return true
	} else {
		fmt.Println("输入不合法")
		return false
	}
}

func (this *Client) updateName() {
	fmt.Println("请输入用户名")
	newName := ""
	fmt.Scanln(&newName)
	if newName == "" {
		fmt.Println("用户名不可为空")
		return
	}
	sendMsg := "rename|" + newName + "\n"
	_, err := this.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write Error", err)
		return
	}

}

func (this *Client) publicChat() {
	fmt.Println("请输入聊天消息，输入exit退出")
	var msg string
	fmt.Scanln(&msg)
	for msg != "exit" {
		if len(msg) != 0 {
			_, err := this.conn.Write([]byte(msg + "\n"))
			if err != nil {
				fmt.Println("conn write error", err)
			}
		}
		msg = ""
		fmt.Println("请输入聊天消息，输入exit退出")
		fmt.Scanln(&msg)
	}
}

func (this *Client) queryUsers() {
	sendMsg := "who\n"
	_, err := this.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn write error", err)
		return
	}
}

func (this *Client) privateChat() {
	this.queryUsers()
	fmt.Println("请输入聊天对象用户名,exit退出")
	user := ""
	fmt.Scanln(&user)
	for user != "exit" {
		fmt.Println("请输入消息,exit退出")
		msg := ""
		fmt.Scanln(&msg)
		for msg != "exit" {
			_, err := this.conn.Write([]byte("to|" + user + "|" + msg + "\n"))
			if err != nil {
				fmt.Println("conn write error", err)
			}
			fmt.Println("请输入消息,exit退出")
			msg = ""
			fmt.Scanln(&msg)
		}
		fmt.Println("请输入聊天对象用户名,exit退出")
		user = ""
		fmt.Scanln(&user)
	}
}

func (this *Client) receiveMsg() {
	io.Copy(os.Stdout, this.conn)
}

func (this *Client) Run() {
	go this.receiveMsg()
	for this.flag != 0 {
		if this.menu() {
			switch this.flag {
			case 0:
				{
					break
				}
			case 1:
				{
					this.publicChat()
					break
				}
			case 2:
				{
					this.privateChat()
					break
				}
			case 3:
				{
					this.updateName()
					break
				}
			}
		}
	}
}
