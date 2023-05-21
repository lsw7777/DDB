package master

import (
	"fmt"
	"log"
	. "Distributed-MiniSQL/common"
)

//在Master服务器中创建表
func (master *Master) CreateTable(args *TableArgs, ip *string) error {
	log.Printf("Master.CreateTable called")
	_, ok := master.tableIP[args.Table]
	if ok {
		log.Printf("%v already exists", args.Table)
		return fmt.Errorf("%v already exists", args.Table)
	}
	bestServer := master.bestServer()
	client := master.regionClients[bestServer]

	var dummy string
	call, err := TimeoutRPC(client.Go("Region.Process", &args.SQL, &dummy, nil), TIMEOUT_S)
	if err != nil {
		log.Printf("%v's Region.Process timeout", bestServer)
		return err 
	}
	if call.Error != nil {
		log.Printf("%v's Region.Process error: %v", bestServer, call.Error)
		return call.Error
	}
	master.addTable(args.Table, bestServer)
	*ip = bestServer
	return nil
}

//在Master服务器中删除表
func (master *Master) DropTable(args *TableArgs, dummyReply *bool) error {
	log.Printf("Master.DropTable called")
	ip, ok := master.tableIP[args.Table]
	if !ok {
		log.Printf("%v not exist", args.Table)
		return fmt.Errorf("%v not exist", args.Table)
	}
	client := master.regionClients[ip]

	var dummy string
	call, err := TimeoutRPC(client.Go("Region.Process", &args.SQL, &dummy, nil), TIMEOUT_S)
	if err != nil {
		log.Printf("%v's Region.Process timeout", ip)
		return err 
	}
	if call.Error != nil {
		log.Printf("%v's Region.Process process error", ip)
		return call.Error 
	}
	master.deleteTable(args.Table, ip)
	return nil
}

//展示Master服务器中所有表
func (master *Master) ShowTables(dummyArgs *bool, tables *[]string) error {
	*tables = make([]string, 0)
	for _, pTables := range master.serverTables {
		*tables = append(*tables, *pTables...)
	}
	return nil
}

//在Master服务器中创建某表的索引
func (master *Master) CreateIndex(args *IndexArgs, ip *string) error {
	log.Printf("Master.CreateIndex called")
	_, ok := master.indexInfo[args.Index]
	if ok {
		log.Printf("%v already exists", args.Index)
		return fmt.Errorf("%v already exists", args.Index)
	}
	*ip = master.tableIP[args.Table]
	client := master.regionClients[*ip]

	var dummy string
	call, err := TimeoutRPC(client.Go("Region.Process", &args.SQL, &dummy, nil), TIMEOUT_S)
	if err != nil {
		log.Printf("%v's Region.Process timeout", *ip)
		return err 
	}
	if call.Error != nil {
		log.Printf("%v's Region.Process error: %v", *ip, call.Error)
		return call.Error 
	}
	master.indexInfo[args.Index] = args.Table
	return nil
}

//在Master服务器中删除某表的索引
func (master *Master) DropIndex(args *IndexArgs, dummyReply *bool) error {
	log.Printf("Master.DropIndex called")
	tbl, ok := master.indexInfo[args.Index]
	if !ok {
		log.Printf("%v not exist", args.Index)
		return fmt.Errorf("%v not exist", args.Index)
	}
	client := master.regionClients[master.tableIP[tbl]]

	var dummy string
	call, err := TimeoutRPC(client.Go("Region.Process", &args.SQL, &dummy, nil), TIMEOUT_S)
	if err != nil {
		log.Printf("%v's Region.Process timeout", master.tableIP[tbl])
		return err 
	}
	if call.Error != nil {
		log.Printf("%v's Region.Process process error", master.tableIP[tbl])
		return call.Error
	}
	delete(master.indexInfo, args.Index)
	return nil
}

//展示Master服务器中所有索引
func (master *Master) ShowIndices(dummyArgs *bool, indices *map[string]string) error {
	log.Printf("Master.ShowIndices called")
	*indices = master.indexInfo
	return nil
}

//获取并保存指定表的服务器IP地址
func (master *Master) TableIP(table *string, ip *string) error {
	log.Printf("Master.TableIP called")
	res, ok := master.tableIP[*table]
	if !ok {
		return fmt.Errorf("%v not exist", *table)
	}
	*ip = res
	return nil
}
