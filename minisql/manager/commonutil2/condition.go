package condition

import (
	"Distributed-MiniSQL/minisql/manager/catalogmanager"
	"fmt"
	"strconv"
	"strings"
)

//定义表行结构体类型
type TableRow struct {
	AttributeValue []string
}

//产生新的表行
func NewTableRow(attributeValue []string) *TableRow {
	return &TableRow{AttributeValue: attributeValue}
}

//增加属性值
func (r *TableRow) AddAttributeValue(attributeValue string) {
	r.AttributeValue = append(r.AttributeValue, attributeValue)
}

//获得属性的索引
func (r *TableRow) GetAttributeValue(index int) string {
	return r.AttributeValue[index]
}

//获得属性大小
func (r *TableRow) GetAttributeSize() int {
	return len(r.AttributeValue)
}


//定义状态结构体类型
type Condition struct {
	Name     string
	Value    string
	Operator string
}

//产生新状态
func NewCondition(name string, operator string, value string) *Condition {
	return &Condition{Name: name, Operator: operator, Value: value}
}

//判断给定的数据行data是否符合Condition对象所描述的条件
func (c *Condition) Satisfy(tableName string, data TableRow) bool {
	index := catalogmanager.GetAttributeIndex(tableName, c.Name)
	type1 := catalogmanager.GetType(tableName, index)
	if type1 == 1 {
		cmpObject := data.GetAttributeValue(index)
		cmpValue := c.Value
		if c.Operator == "=" {
			flag := cmpObject == cmpValue
			return flag
		} else if c.Operator == "<>" {
			return strings.Compare(cmpObject, cmpValue) != 0
		} else if c.Operator == ">" {
			return strings.Compare(cmpObject, cmpValue) > 0
		} else if c.Operator == "<" {
			return strings.Compare(cmpObject, cmpValue) < 0
		} else if c.Operator == ">=" {
			return strings.Compare(cmpObject, cmpValue) >= 0
		} else if c.Operator == "<=" {
			return strings.Compare(cmpObject, cmpValue) <= 0
		} else {
			return false
		}
	} else if type1 == 3 {
		cmpObject, err := strconv.ParseFloat(data.GetAttributeValue(index), 64)
		if err != nil {
			fmt.Print(err)
			return false
		}
		cmpValue, err := strconv.ParseFloat(c.Value, 64)
		if err != nil {
			fmt.Print(err)
			return false
		}
		if c.Operator == "=" {
			return cmpObject == cmpValue
		} else if c.Operator == "<>" {
			return cmpObject != cmpValue
		} else if c.Operator == ">" {
			return cmpObject > cmpValue
		} else if c.Operator == "<" {
			return cmpObject < cmpValue
		} else if c.Operator == ">=" {
			return cmpObject >= cmpValue
		} else if c.Operator == "<=" {
			return cmpObject <= cmpValue
		} else {
			return false
		}
	} else if type1 == 2 {
		cmpObject, err := strconv.Atoi(data.GetAttributeValue(index))
		if err != nil {
			fmt.Print(err)
			return false
		}
		cmpValue, err := strconv.Atoi(c.Value)
		if err != nil {
			fmt.Print(err)
			return false
		}
		switch c.Operator {
		case "=":
			return cmpObject == cmpValue
		case "<>":
			return cmpObject != cmpValue
		case ">":
			return cmpObject > cmpValue
		case "<":
			return cmpObject < cmpValue
		case "<=":
			return cmpObject <= cmpValue
		case ">=":
			return cmpObject >= cmpValue
		default:
			return false
		}
	} else {
		return false
	}
}
