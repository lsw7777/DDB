package catalogmanager

import( 
	index "Distributed-MiniSQL/minisql/manager/commonutil"
)

//定义数据类型结构体类型
type Datatype struct {
	NumType int // 1为char,2 为int, 3为float
	length  int //长度
}

//定义新的数据类型
func NewDatatype(NumType int, length int) *Datatype { 
	if NumType == 1 {
		return &Datatype{NumType: 1, length: length}
	}
	if NumType == 2 {
		return &Datatype{NumType: 2, length: 4}
	} else {
		return &Datatype{NumType: 3, length: 4}
	}
}



//定义属性结构体类型
type Attribute struct {
	AttributeName string
	FieldType     Datatype
	IsUnique      bool //该属性是否是unique的
}

//定义新属性
func NewAttribute(attributeName string, NumType int, length int, isUnique bool) *Attribute {
	newDatatype := NewDatatype(NumType, length)
	return &Attribute{AttributeName: attributeName, FieldType: *newDatatype, IsUnique: isUnique}
}

//定义地址结构体类型
type Address struct {
	FileName    string
	BlockOffset int //file中的block offset
	ByteOffset  int //block 中的
}

//定义新地址
func NewAddress(fileName string, blockOffset int, byteOffset int) *Address {
	return &Address{FileName: fileName, BlockOffset: blockOffset, ByteOffset: byteOffset}
}

//用于比较两个地址
func (a Address) Compare(b Address) bool {
	if a.FileName == b.FileName {
		if a.BlockOffset == b.BlockOffset {
			return a.ByteOffset > b.ByteOffset
		} else {
			return a.BlockOffset > b.BlockOffset
		}
	} else {
		return a.FileName > b.FileName
	}
}

//定义表结构体类型
type Table struct {
	TableName       string
	PrimaryKey      string
	AttributeVector []Attribute
	IndexVector     []index.Index
	IndexNum        int
	AttributeNum    int
	RowNum          int
	RowLength       int
}

//定义新表
func NewTable(tableName string, primaryKey string, attributeVector []Attribute) *Table {
	var table = Table{TableName: tableName, PrimaryKey: primaryKey,
		AttributeVector: attributeVector, AttributeNum: len(attributeVector)}
	for i := 0; i < len(table.AttributeVector); i++ {
		if table.AttributeVector[i].AttributeName == primaryKey {
			table.AttributeVector[i].IsUnique = true
		}
		table.RowLength += table.AttributeVector[i].FieldType.length //计算一行数据所需的存储空间
	}
	return &table
}

//也是定义新表。新增两个参数：索引向量和行数。
func NewTable2(tableName string, primaryKey string, attributeVector []Attribute, indexVector []index.Index, rowNum int) *Table {
	var table = Table{TableName: tableName, PrimaryKey: primaryKey,
		AttributeVector: attributeVector, AttributeNum: len(attributeVector), IndexVector: indexVector, IndexNum: len(indexVector), RowNum: rowNum}
	for i := 0; i < len(table.AttributeVector); i++ {
		if table.AttributeVector[i].AttributeName == primaryKey {
			table.AttributeVector[i].IsUnique = true
		}
		table.RowLength += table.AttributeVector[i].FieldType.length //计算一行数据所需的存储空间
	}
	return &table
}
