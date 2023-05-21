package main

import (
	. "Distributed-MiniSQL/client"
	"os"
)

//主函数
func main() {
	var client Client
	client.Init(os.Args[1])
	client.Run()
}
