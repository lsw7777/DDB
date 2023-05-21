package main

import (
	. "Distributed-MiniSQL/client"
	"os"
)

//主函数，开启client
func main() {
	var client Client
	client.Init(os.Args[1])
	client.Run()
}
