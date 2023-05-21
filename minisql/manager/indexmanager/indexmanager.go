package indexmanager

import (
	"Distributed-MiniSQL/common"
	buffermanager "Distributed-MiniSQL/minisql/manager/buffermanager"
	catalogmanager "Distributed-MiniSQL/minisql/manager/catalogmanager"
	index "Distributed-MiniSQL/minisql/manager/commonutil"
	condition "Distributed-MiniSQL/minisql/manager/commonutil2"
	"bufio"
	"fmt"
	"os"
	"strconv"
)


var TreeMap = make(map[string]BPTree)

func InitIndex() {
	FileName := common.DIR + "indexCatalog.txt"
	file, err := os.Open(FileName)
	if err != nil {
		fmt.Println("文件打开失败")
		return
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	for i := 0; i < len(lines); {
		tempIndexName := lines[i]
		tempTableName := lines[i+1]
		tempAttributeName := lines[i+2]
		tempBlockNum, _ := strconv.Atoi(lines[i+3])
		tempRootNum, _ := strconv.Atoi(lines[i+4])
		CreateIndex(*index.NewIndex2(tempIndexName, tempTableName, tempAttributeName, tempBlockNum, tempRootNum))
		i += 5
	}
	fmt.Println("索引信息加载成功")
}

func CreateIndex(idx index.Index) bool {
	fileName := common.DIR + idx.IndexName + ".index" 
	BuildIndex(idx)
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("打开文件错误")
		fmt.Println(err)
		return false
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	fmt.Fprintln(w, idx.IndexName)
	fmt.Fprintln(w, idx.TableName)
	fmt.Fprintln(w, idx.AttributeName)
	fmt.Fprintln(w, strconv.Itoa(idx.BlockNum))
	fmt.Fprint(w, strconv.Itoa(idx.RootNum))
	w.Flush()
	return true
}

func DropIndex(idx index.Index) bool {
	fileName := common.DIR + idx.IndexName + ".index"
	err := os.Remove(fileName) //删除文件
	if err != nil {
		fmt.Println("删除失败")
		return false
	}
	delete(TreeMap, idx.IndexName)
	return true
}

func BuildIndex(idx index.Index) {
	tableName := idx.TableName
	attributeName := idx.AttributeName
	tupleNum := catalogmanager.GetRowNum(tableName)
	tupleLength := GetStoreLength(tableName)
	blockOffset := 0
	processNum := 0
	byteOffset := 4
	IndexNum := catalogmanager.GetAttributeIndex(tableName, attributeName)

	tree := NewBPTree()
	// dataType := catalogmanager.GetType(tableName, IndexNum)
	block := buffermanager.ReadBlockFromDiskQuote(tableName, 0)
	for processNum < tupleNum {
		if byteOffset+tupleLength >= buffermanager.BLOCKSIZE { 
			blockOffset++
			byteOffset = 0 
			block = buffermanager.ReadBlockFromDiskQuote(tableName, blockOffset)
			if block == nil {
				
			}
		}
		if block.ReadInteger(byteOffset) < 0 {
			value := catalogmanager.NewAddress(tableName, blockOffset, byteOffset)
			row := GetTuple(tableName, *block, byteOffset)
			key := row[IndexNum]

			tree.Insert(key, *value)
			processNum++
		}
		byteOffset += tupleLength
	}
	TreeMap[idx.IndexName] = *tree
}

//需要抛出错误
func Select(idx index.Index, cond condition.Condition) []catalogmanager.Address {
	tree := TreeMap[idx.IndexName]
	indexNum := catalogmanager.GetAttributeIndex(idx.TableName, idx.AttributeName)
	Datatype := catalogmanager.GetType(idx.TableName, indexNum) //类型
	operator := cond.Operator
	var value interface{}
	if Datatype == 1 { //char
		value = cond.Value
	} else if Datatype == 2 { //int
		value, _ = strconv.Atoi(cond.Value)
	} else if Datatype == 3 { //float
		value, _ = strconv.ParseFloat(cond.Value, 64)
	}
	if operator == "=" {
		result := tree.FindEq(value)
		if result == nil {
			return nil
		}
		return []catalogmanager.Address{*result}
	} else if operator == "<>" {
		return tree.FindNeq(value)
	} else if operator == ">" {
		return tree.FindGreater(value)
	} else if operator == "<" {
		return tree.FindLess(value)
	} else if operator == ">=" {
		return tree.FindGeq(value)
	} else if operator == "<=" {
		return tree.FindLeq(value)
	} else {
		
	}
	return []catalogmanager.Address{}
}

func Delete(idx index.Index, key string) {
	tree := TreeMap[idx.IndexName]
	indexNum := catalogmanager.GetAttributeIndex(idx.TableName, idx.AttributeName)
	Datatype := catalogmanager.GetType(idx.TableName, indexNum) 
	var value interface{}
	if Datatype == 1 { //char
		value = key
	} else if Datatype == 2 { //int
		value, _ = strconv.Atoi(key)
	} else if Datatype == 3 { //float
		value, _ = strconv.ParseFloat(key, 64)
	}
	tree.Delete(value)
}

func Insert(idx index.Index, key string, value catalogmanager.Address) {
	tree := TreeMap[idx.IndexName]
	indexNum := catalogmanager.GetAttributeIndex(idx.TableName, idx.AttributeName)
	Datatype := catalogmanager.GetType(idx.TableName, indexNum) 
	var value1 interface{}
	if Datatype == 1 { //char
		value1 = key
	} else if Datatype == 2 { //int
		value1, _ = strconv.Atoi(key)
	} else if Datatype == 3 { //float
		value1, _ = strconv.ParseFloat(key, 64)
	}
	tree.Insert(value1, value)
}

func Update(idx index.Index, key string, value catalogmanager.Address) {
	tree := TreeMap[idx.IndexName]
	indexNum := catalogmanager.GetAttributeIndex(idx.TableName, idx.AttributeName)
	Datatype := catalogmanager.GetType(idx.TableName, indexNum) 
	var value1 interface{}
	if Datatype == 1 { //char
		value1 = key
	} else if Datatype == 2 { //int
		value1, _ = strconv.Atoi(key)
	} else if Datatype == 3 { //float
		value1, _ = strconv.ParseFloat(key, 64)
	}
	tree.Insert(value1, value)
}


func GetStoreLength(tableName string) int {
	rowLen := catalogmanager.GetRowLength(tableName)
	if rowLen > 4 { // a valid byte
		return rowLen + 1
	} else {
		return 5 //4 + 1
	}
}

func GetTuple(tableName string, block buffermanager.Block, offset int) []interface{} {
	attributeNum := catalogmanager.GetAttributeNum(tableName)
	var attributeValue interface{}
	var result []interface{}

	offset++ 
	for i := 0; i < attributeNum; i++ {
		length := catalogmanager.GetLength2(tableName, i)
		datatype := catalogmanager.GetType(tableName, i)
		if datatype == 1 { 
			attributeValue = block.ReadString(offset, length)

			attributeValue = rmu0000(attributeValue.(string))
		} else if datatype == 2 { 
			attributeValue = int(block.ReadInteger(offset))
		} else if datatype == 3 { 
			attributeValue = float32(block.ReadFloat(offset))
		}
		offset += length
		result = append(result, attributeValue)
	}
	return result
}

func rmu0000(s string) string {
	str := make([]rune, 0, len(s))
	for _, v := range []rune(s) {
		if v == 0 {
			continue
		}
		str = append(str, v)
	}
	return string(str)
}
