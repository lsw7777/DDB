package region

import (
	"fmt"
	"log"
	"net/rpc"

	. "Distributed-MiniSQL/common"
	"Distributed-MiniSQL/minisql/manager/api"
	"Distributed-MiniSQL/minisql/manager/interpreter"
)


//用于处理接收到的 SQL 语句，它会调用 region.processSQL 方法来将 SQL 语句转换为对 etcd 的操作，并返回执行结果。如果执行出错，则会返回错误信息
func (region *Region) Process(input *string, reply *string) error {
	log.Printf("Region.Process called: %v", *input)
	res, err := region.processSQL(*input)
	if err != nil {
		log.Printf("Region.Process failed")
		return fmt.Errorf("Region.Process failed")
	} else {
		*reply = res
		if region.backupIP != "" {
			rpcBackupRegion, err := rpc.DialHTTP("tcp", region.backupIP+REGION_PORT)
			if err != nil {
				log.Printf("fail to connect to backup %v", region.backupIP)
				return nil
			}
			// backup's Region.Process must return nil error
			_, err = TimeoutRPC(rpcBackupRegion.Go("Region.Process", &input, &reply, nil), TIMEOUT_S)
			if err != nil {
				log.Printf("%v's Region.Process timeout", region.backupIP)
				return nil
			}
		}
	}
	return err
}

//用于将当前节点设为备份节点，并连接到指定的备份节点
func (region *Region) AssignBackup(ip *string, dummyReply *bool) error {
	log.Printf("Region.AssignBackup called: backup ip: %v", *ip)
	client, err := rpc.DialHTTP("tcp", *ip+REGION_PORT)
	if err != nil {
		log.Printf("rpc.DialHTTP err: %v", *ip+REGION_PORT)
	} else {
		region.backupClient = client
		region.backupIP = *ip
		_, err = TimeoutRPC(region.backupClient.Go("Region.DownloadSnapshot", &region.hostIP, &dummyReply, nil), TIMEOUT_L)
		if err != nil {
			log.Printf("%v's Region.DownloadSnapshot timeout", *ip)
			region.RemoveBackup(nil, nil)
		}
	}
	return err
}

//用于将备份IP地址从当前Region结构体中移除，并关闭与备份客户端的连接。
func (region *Region) RemoveBackup(dummyArgs, dummyReply *bool) error {
	log.Printf("Region.RemoveBackup called: remove %v", region.backupIP)
	region.backupIP = ""
	if region.backupClient != nil {
		region.backupClient.Close()
	}
	region.backupClient = nil
	return nil
}

//用于下载指定IP的快照，并清除Region结构体中的备份IP。
func (region *Region) DownloadSnapshot(ip *string, dummyReply *bool) error {
	log.Printf("Region.DownloadSnapshot called: download %v's snapshot", *ip)
	region.RemoveBackup(nil, nil)
	region.fu.DownloadDir(WORKING_DIR+DIR, DIR, *ip)
	api.Initial()
	return nil
}

//用于将接收到的SQL查询语句作为参数，并将其解析和处理为结果字符串
func (region *Region) processSQL(sql string) (string, error) {
	res := interpreter.Interpret(sql)

	if res == "" {
		return res, fmt.Errorf("process failed")
	}

	api.Store()
	return res, nil
}
