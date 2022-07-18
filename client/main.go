package main

import (
	"flag"
	"fmt"
)

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "h", "127.0.0.1", "设置服务器IP地址(默认127.0.0.1)")
	flag.IntVar(&serverPort, "p", 8088, "设置服务器端口(默认8088)")
}

func main() {
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>连接服务器失败")
		return
	}
	fmt.Println(">>>>>连接服务器成功")
	client.Run()
}
