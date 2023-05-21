package main

import (
	. "Distributed-MiniSQL/master"
	"os"
	"strconv"
)

//主函数，开启master
func main() {
	var master Master
	regionCount, _ := strconv.ParseInt(os.Args[1], 10, 0)
	master.Init(int(regionCount))
	master.Run()
}
