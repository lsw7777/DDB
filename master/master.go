package master

import (
	"math"
	"net"
	"net/http"
	"net/rpc"
	"time"
	clientv3 "go.etcd.io/etcd/client/v3"
	. "Distributed-MiniSQL/common"
)


//定义Master结构体变量
type Master struct {
	regionCount int
	etcdClient    *clientv3.Client
	regionClients map[string]*rpc.Client
	serverTables map[string]*[]string 
	tableIP      map[string]string   
	indexInfo    map[string]string    
	backupInfo map[string]string
}

//初始化并返回一个Master变量
func (master *Master) Init(regionCount int) {
	master.regionCount = regionCount

	master.regionClients = make(map[string]*rpc.Client)
	master.serverTables = make(map[string]*[]string)
	master.tableIP = make(map[string]string)
	master.indexInfo = make(map[string]string)

	master.backupInfo = make(map[string]string)
}

//启动Master主节点对外提供RPC服务，不断等待新的请求的到来
func (master *Master) Run() {
	// connect to local etcd server
	master.etcdClient, _ = clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://" + HOST_ADDR},
		DialTimeout: 1 * time.Second,
	})
	defer master.etcdClient.Close()
	go master.watch()

	rpc.Register(master)
	rpc.HandleHTTP()
	l, _ := net.Listen("tcp", MASTER_PORT)
	go http.Serve(l, nil)

	for {
		time.Sleep(10 * time.Second)
	}
}

//在Master服务器中增加一个表
func (master *Master) addTable(table, ip string) {
	master.tableIP[table] = ip
	AddUniqueToSlice(master.serverTables[ip], table)
}

//在Master服务器中删除一个表
func (master *Master) deleteTable(table, ip string) {
	master.deleteTableIndices(table)
	delete(master.tableIP, table)
	DeleteFromSlice(master.serverTables[ip], table)
}

//删除 Master服务器中某个表的所有索引
func (master *Master) deleteTableIndices(table string) {
	targets := make([]string, 0)
	for idx, tbl := range master.indexInfo {
		if tbl == table {
			targets = append(targets, idx)
		}
	}
	for _, idx := range targets {
		delete(master.indexInfo, idx)
	}
}

//选择当前可用且任务负载最轻的服务器，并返回其 IP 地址
func (master *Master) bestServer() string {
	min, res := math.MaxInt, ""
	for ip, pTables := range master.serverTables {
		if len(*pTables) < min {
			min, res = len(*pTables), ip
		}
	}
	return res
}

//将一个服务器上的所有表格转移到另一个服务器，并在目的服务器中更新相关信息
func (master *Master) transferServerTables(src, dst string) {
	pTables := master.serverTables[src]
	for _, table := range *pTables {
		master.tableIP[table] = dst
	}
	master.serverTables[dst] = pTables
	delete(master.serverTables, src)
}

//删除指定IP地址的服务器上的所有表
func (master *Master) removeServerTables(ip string) {
	pTables := master.serverTables[ip]
	for _, table := range *pTables {
		master.deleteTableIndices(table)
		delete(master.tableIP, table)
	}
	delete(master.serverTables, ip)
}
