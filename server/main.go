package main

func main() {
	server := BuildServer("127.0.0.1", 8088)
	server.start()
}
