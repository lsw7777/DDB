package region

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"
	clientv3 "go.etcd.io/etcd/client/v3"
	. "Distributed-MiniSQL/common"
	"Distributed-MiniSQL/minisql/manager/api"
)

//定义Region结构体变量
type Region struct {
	masterIP string				//master服务器的IP地址
	hostIP   string				//region服务器的IP地址
	backupIP string				//该region服务器的从服务器的IP地址
	etcdClient   *clientv3.Client		//媒介
	masterClient *rpc.Client		//媒介
	backupClient *rpc.Client		//媒介
	fu           FtpUtils
}

//用于初始化 Region 结构体
func (region *Region) Init(masterIP, hostIP string) {
	region.masterIP = masterIP
	region.hostIP = hostIP

	region.fu.Construct()

	api.Initial()
}

//开启了一个stayOnline()任务，还开启了一个RPC服务任务
func (region *Region) Run() {
	// connect to local etcd server
	region.etcdClient, _ = clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://" + HOST_ADDR},
		DialTimeout: 5 * time.Second,
	})
	defer region.etcdClient.Close()
	go region.stayOnline()
	rpc.Register(region)
	rpc.HandleHTTP()
	l, _ := net.Listen("tcp", REGION_PORT)
	go http.Serve(l, nil)
	client, err := rpc.DialHTTP("tcp", region.masterIP+MASTER_PORT)
	if err != nil {
		log.Printf("rpc.DialHTTP err: %v", region.masterIP+MASTER_PORT)
		return
	}
	region.masterClient = client

	for {
		time.Sleep(10 * time.Second)
	}
}

//是一个死循环，通过 etcd 的 lease 功能和 keepAlive 功能来实现在线状态的维护，防止节点宕机后 etcd 上的信息仍残留
func (region *Region) stayOnline() {
	for {
		resp, err := region.etcdClient.Grant(context.Background(), 5)
		if err != nil {
			log.Printf("etcd grant error")
			continue
		}

		_, err = region.etcdClient.Put(context.Background(), region.hostIP, "", clientv3.WithLease(resp.ID))
		if err != nil {
			log.Printf("etcd put error")
			continue
		}

		ch, err := region.etcdClient.KeepAlive(context.Background(), resp.ID)
		if err != nil {
			log.Printf("etcd keepalive error")
			continue
		}

		for _ = range ch {
		}
	}
}
