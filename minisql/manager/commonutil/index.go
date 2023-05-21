package index

//定义索引结构体类型
type Index struct {
	IndexName     string
	TableName     string
	AttributeName string

	RootNum  int
	BlockNum int
}

//产生新索引
func NewIndex(indexName string, tableName string, attributeName string) *Index {
	return &Index{IndexName: indexName, TableName: tableName, AttributeName: attributeName}
}

//产生新索引
func NewIndex2(indexName string, tableName string, attributeName string, blockNum int, rootNum int) *Index {
	return &Index{IndexName: indexName, TableName: tableName, AttributeName: attributeName, BlockNum: blockNum, RootNum: rootNum}
}
