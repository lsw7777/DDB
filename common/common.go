package common

import (
	"fmt"
	"io/ioutil"
	"net/rpc"
	"os"
	"time"
)

type Identity int

//常量定义
const (
	//超时
	TIMEOUT_S = 1000
	TIMEOUT_M = 2000
	TIMEOUT_L = 10000

	//网络连接
	NETWORK = "tcp"

	MASTER_PORT = ":4733"
	REGION_PORT = ":2016"

	//Etcd的主机地址
	HOST_ADDR = "127.0.0.1:2379"

	//工作区域
	WORKING_DIR = "distributed-mini-sql/"
	DIR         = "sql/"
)

//定义"表命令"结构体类型
type TableArgs struct {
	Table string
	SQL   string
}

//定义"索引表命令"结构体类型
type IndexArgs struct {
	Index string
	Table string
	SQL   string
}

//在字符串数组中查找指定字符串的序号值
func FindElement(pSlice *[]string, str string) int {
	for i, elem := range *pSlice {
		if elem == str {
			return i
		}
	}
	return -1
}


//向字符串数组中添加指定字符串（前提是字符串数组中没有该指定字符串）
func AddUniqueToSlice(pSlice *[]string, str string) {
	if FindElement(pSlice, str) == -1 {
		*pSlice = append(*pSlice, str)
	}
}


//从字符串数组中删除指定字符串
func DeleteFromSlice(pSlice *[]string, str string) bool {
	index := FindElement(pSlice, str)
	if index == -1 {
		return false
	}
	(*pSlice)[index] = (*pSlice)[len(*pSlice)-1]
	*pSlice = (*pSlice)[:len(*pSlice)-1]
	return true
}

//返回自身的IP地址
func GetHostIP() string {
	return ""
}

//删除指定文件名的本地文件
func DeleteLocalFile(fileName string) {
	err := os.Remove(fileName)
	if err != nil {
		fmt.Printf("delete local file failed: %v\n", err)
	}
}

//非阻塞的RPC超时处理函数
func TimeoutRPC(call *rpc.Call, ms int) (*rpc.Call, error) {
	select {
	case res := <-call.Done:
		return res, nil
	case <-time.After(time.Duration(ms) * time.Millisecond):
		return nil, fmt.Errorf("%v timeout", call.ServiceMethod)
	}
}

//清理目录
func CleanDir(localDir string) {
	dir, err := ioutil.ReadDir(localDir)
	if err != nil {
		fmt.Println("Can't obtain files in dir")
	}
	for _, d := range dir {
		os.RemoveAll(localDir + d.Name())
	}
}
