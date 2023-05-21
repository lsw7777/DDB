package client

import (
	. "Distributed-MiniSQL/common"
	"bufio"
	"errors"
	"fmt"
	"net/rpc"
	"os"
	"strings"
)

//取代字符串
func ReplaceAll(str, old, new string) string {
	return strings.Replace(str, old, new, -1)
}

//客户端结构体
type Client struct {
	ipCache      map[string]string
	rpcMaster    *rpc.Client
	rpcRegionMap map[string]*rpc.Client // [ip]rpc
}

type TableOp int

//定义SQL不同操作的OPCODE值
const (
	CREATE_TBL = 0
	DROP_TBL   = 1
	SHOW_TBL   = 2
	CREATE_IDX = 3
	DROP_IDX   = 4
	SHOW_IDX   = 5
	OTHERS     = 6
)

//初始化并返回一个client结构体变量
func (client *Client) Init(masterIP string) {
	client.ipCache = make(map[string]string)
	client.rpcRegionMap = make(map[string]*rpc.Client)

	rpcMas, err := rpc.DialHTTP("tcp", masterIP+MASTER_PORT)
	if err != nil {
		fmt.Printf("SYSTEM HINT>>> client connect error: %v", err)
	}
	client.rpcMaster = rpcMas

}

//根据预处理得到的信息，实现用户SQL输入命令与数据库的交互
func (client *Client) Run() {
	execFileMode := false
	var commands []string
	indexOfCommand := 0
	for {
		input := ""
		if execFileMode {
			input = commands[indexOfCommand]
			input = strings.Trim(input, " ")
			indexOfCommand += 1
			if indexOfCommand == len(commands) { // come to the final line, return to input by keyboard mode
				indexOfCommand = 0
				execFileMode = false
				commands = commands[0:0]
			}
		} else {
			fmt.Println("new message>>> input your SQL command: ")
			part_input := "" // part of the input, all of them compose input

			for len(part_input) == 0 || part_input[len(part_input)-1] != ';' {
				part_input, _ = bufio.NewReader(os.Stdin).ReadString('\n')
				part_input = strings.TrimRight(part_input, "\r\n")
				if len(part_input) == 0 {
					continue
				}
				input += part_input
				input += " "
			}

			input = strings.Trim(input, " ")

			if strings.HasPrefix(input, "quit") {
				quitCommand := strings.Trim(input, "; ")
				if quitCommand == "quit" {
					// client exit directly
					fmt.Println("new message>>> You choose to quit, bye!")
					break
				} else {
					fmt.Println("SYSTEM HINT>>> invalid format of SQL command")
					continue
				}
			} else {
				command := ReplaceAll(input, "\\s+", " ")
				words := strings.Split(command, " ")
				if len(words) == 2 && words[0] == "execfile" {
					fileName := words[1]
					fileName = fileName[0:strings.Index(fileName, ";")]
					f, err := os.Open(fileName)
					if err != nil {
						fmt.Println("INPUT FORMAT ERROR>>> choose execfile but the file can't be found")
					} else {
						scanner := bufio.NewScanner(f)
						for scanner.Scan() {
							commands = append(commands, scanner.Text())
						}
						execFileMode = true
						indexOfCommand = 0
						fmt.Println("HINT>>> start to execute command in file")
					}
					continue // no matter execfile or or command error, jump over this turn
				}
			}
		}

		op, table, index, err := client.preprocessInput(input)
		if err != nil {
			fmt.Printf("INPUT FORMAT ERROR>>> %v\n", err)
			continue
		}
		switch op {
		case CREATE_TBL:
			// call Master.CreateTable rpc
			args, ip := TableArgs{Table: table, SQL: input}, ""
			call, err := TimeoutRPC(client.rpcMaster.Go("Master.CreateTable", &args, &ip, nil), TIMEOUT_M)
			if err != nil {
				fmt.Println("SYSTEM HINT>>> timeout, master down!")
				break
			}
			if call.Error != nil {
				fmt.Println("RESULT>>> create table failed")
			} else {
				fmt.Println("RESULT>>> create table succeed, table in ip: " + ip)
			}
		case DROP_TBL:
			args, dummy := TableArgs{Table: table, SQL: input}, false
			call, err := TimeoutRPC(client.rpcMaster.Go("Master.DropTable", &args, &dummy, nil), TIMEOUT_M)
			if err != nil {
				fmt.Println("SYSTEM HINT>>> timeout, master down!")
				break
			}
			if call.Error != nil {
				fmt.Println("RESULT>>> drop table failed")
			} else {
				fmt.Println("RESULT>>> drop table succeed")
			}
		case SHOW_TBL:
			var dummyArgs bool
			var tables []string
			// tables = make([]string, 0)
			call, err := TimeoutRPC(client.rpcMaster.Go("Master.ShowTables", &dummyArgs, &tables, nil), TIMEOUT_S)
			if err != nil {
				fmt.Println("SYSTEM HINT>>> timeout, master down!")
				break
			}
			if call.Error != nil {
				fmt.Println("RESULT>>> show tables failed")
			} else {
				fmt.Println("RESULT>>> tables in region:\n---------------")
				for _, table := range tables {
					fmt.Printf("|  %v  |\n", table)
				}
				fmt.Println("---------------")
			}

		case CREATE_IDX:
			args, ip := IndexArgs{Index: index, Table: table, SQL: input}, ""
			call, err := TimeoutRPC(client.rpcMaster.Go("Master.CreateIndex", &args, &ip, nil), TIMEOUT_M)
			if err != nil {
				fmt.Println("SYSTEM HINT>>> timeout, master down!")
				break
			}
			if call.Error != nil {
				fmt.Println("RESULT>>> create indexes failed")
			} else {
				fmt.Printf("RESULT>>> create index succeed on table %v in %v\n", table, ip)
			}
		case DROP_IDX:
			args, dummy := IndexArgs{Index: index, SQL: input}, false
			call, err := TimeoutRPC(client.rpcMaster.Go("Master.DropIndex", &args, &dummy, nil), TIMEOUT_M)
			if err != nil {
				fmt.Println("SYSTEM HINT>>> timeout, master down!")
				break
			}
			if call.Error != nil {
				fmt.Println("RESULT>>> drop indexes failed")
			} else {
				fmt.Println("RESULT>>> drop indexes succeed")
			}
		case SHOW_IDX:
			var dummyArgs bool
			var indices map[string]string
			call, err := TimeoutRPC(client.rpcMaster.Go("Master.ShowIndices", &dummyArgs, &indices, nil), TIMEOUT_S)
			if err != nil {
				fmt.Println("SYSTEM HINT>>> timeout, master down!")
				break
			}
			if call.Error != nil {
				fmt.Println("RESULT>>> show indices failed")
			} else {
				fmt.Println("RESULT>>> indices in region:\n|    index    |    table    |")
				for index, table := range indices {
					fmt.Printf("|   %v   |   %v   |\n", index, table)
				}
				fmt.Println("---------------")
			}
		case OTHERS:
			ip, ok := client.ipCache[table]
			if !ok {
				ip = client.updateCache(table)
				if ip == "" {
					// fmt.Println("can't find the corresponding ip in cache")
					break
				}
			} else {
				fmt.Println("SYSTEM HINT>>> first phase: find corresponding ip in table-ip cache: " + ip)
			}
			result := ""
			rpcRegion, ok := client.rpcRegionMap[ip]
			if !ok {
				rpcRegion, err = rpc.DialHTTP("tcp", ip+REGION_PORT)
				if err != nil {
					delete(client.ipCache, table)
					break
				} else {
					client.rpcRegionMap[ip] = rpcRegion
				}
			}

			call, err := TimeoutRPC(rpcRegion.Go("Region.Process", &input, &result, nil), TIMEOUT_S)
			if err != nil || call.Error != nil {
				if err != nil {
					fmt.Println("SYSTEM HINT>>> region process timeout")
					fmt.Println("SYSTEM HINT>>> start to reconnect with second phase operation...")
				} else {
					fmt.Println("SYSTEM HINT>>> can't obtain result in first phase, maybe old ip or SQL command error")
					fmt.Println("SYSTEM HINT>>> start second phase operation...")
				}

				delete(client.ipCache, table)
				delete(client.rpcRegionMap, ip)
				new_ip := client.updateCache(table) // obtain newest ip
				if new_ip == "" {
					break
				}
				new_rpcRegion, err := rpc.DialHTTP("tcp", new_ip+REGION_PORT)
				if err != nil {
					fmt.Printf("SYSTEM HINT>>> fail to connect to region: " + ip)
					delete(client.ipCache, table)
					break
				}
				call, err := TimeoutRPC(new_rpcRegion.Go("Region.Process", &input, &result, nil), TIMEOUT_S)
				if err != nil {
					fmt.Println("SYSTEM HINT>>> region process timeout")
					break
				}
				if call.Error != nil {
					fmt.Println("SYSTEM HINT>>> can't obatin result, maybe input is error")
					break
				}
				client.ipCache[table] = new_ip
				client.rpcRegionMap[new_ip] = new_rpcRegion
				fmt.Println("SYSTEM HINT>>> update ip: " + ip + " and add " + new_ip + " to iptablemap")
			}
			if result != "" {
				fmt.Println("RESULT>>> \n" + result)
			}
		}

	}
}

//预处理输入的SQL语句，提取出操作类型、表名、索引名、可能发生的错误类型并返回
func (client *Client) preprocessInput(input string) (TableOp, string, string, error) {
	input = strings.Trim(input, ";")
	table := ""
	index := ""
	var op TableOp
	op = OTHERS
	var err error
	err = nil
	input = ReplaceAll(input, "\\s+", " ")
	words := strings.Split(input, " ")
	if words[0] == "create" {
		if len(words) < 3 {
			err = errors.New("lack of input in create operation")
			return op, table, index, err
		}

		if words[1] == "table" {
			op = CREATE_TBL
			table = words[2]
			if strings.Contains(table, "(") { 
				table = table[0:strings.Index(table, "(")]
			}
		} else if words[1] == "index" {
			op = CREATE_IDX
			if len(words) < 5 {
				err = errors.New("lack of words in create index")
				return op, table, index, err
			}
			index = words[2]
			table = words[4]
			if strings.Contains(table, "(") { 
				table = table[0:strings.Index(table, "(")]
			}
		} else {
			err = errors.New("the type you create can't be recognized")
		}
		return op, table, index, err
	} else if words[0] == "drop" {
		if len(words) == 3 {
			if words[1] == "table" {
				op = DROP_TBL
				table = words[2]
			} else if words[1] == "index" {
				op = DROP_IDX
				index = words[2]
			} else {
				err = errors.New("please drop table or index")
			}
		} else {
			err = errors.New("number of words false in drop")
		}
		return op, table, index, err
	} else if words[0] == "show" {
		if len(words) == 2 {
			if words[1] == "tables" {
				op = SHOW_TBL
			} else if words[1] == "indexes" {
				op = SHOW_IDX
			} else {
				err = errors.New("show command only specify indexes/tables")
			}
			return op, table, index, err
		} else {
			err = errors.New("show command only need to specify indexes/tables")
		}
	} else {
		op = OTHERS
		if words[0] == "select" {
			//select语句的表名放在from后面
			for i := 0; i < len(words); i++ {
				if words[i] == "from" && i != (len(words)-1) {
					table = words[i+1]
					break
				}
			}
		} else if words[0] == "insert" || words[0] == "delete" {
			if len(words) >= 3 {
				table = words[2]
			}
		} else {
			err = errors.New("can't recoginize your command")
		}
	}

	// 只要table仍为""，说明没拿到表名
	if table == "" && op == OTHERS && err == nil {
		err = errors.New("no table name in input")
	}
	return op, table, index, err
}

//这里目前还没有考虑没有查到ip的情况
func (client *Client) updateCache(table string) string {
	ip := ""
	// call Master.TableIP rpc
	call, err := TimeoutRPC(client.rpcMaster.Go("Master.TableIP", &table, &ip, nil), TIMEOUT_S)
	if err != nil {
		fmt.Println("SYSTEM HINT>>> during update cache, timeout, master down")
		return ip
	}
	if call.Error != nil {
		fmt.Println("SYSTEM HINT>>> table invalid, can't update cache")
		return ip
	}
	client.ipCache[table] = ip
	return ip
}
