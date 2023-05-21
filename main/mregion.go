package main

import (
	. "Distributed-MiniSQL/region"
	"os"
)

//主函数，开启region
func main() {
	var region Region
	region.Init(os.Args[1], os.Args[2])
	region.Run()
}
