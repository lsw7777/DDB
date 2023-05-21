package qexception

import (
	"fmt"
	"strconv"
)

var Ex = []string{"Syntax error", "Run time error"}

//定义异常结构体类型
type Qexception struct {
	DataType int    //exception type: 0 for 'syntax error' and 1 for 'rn time error'
	Status   int    //status code
	Msg      string //exception message
}

//新建异常
func newqexception(dataType int, status int, msg string) *Qexception {
	return &Qexception{
		Status:   status,
		DataType: If(dataType >= 0 && dataType <= len(Ex), dataType, 0),
		Msg:      msg,
	}
}

//模拟if语句，是则返回trueVal，否则返回falseVal
func If(condition bool, trueVal int, falseVal int) int {
	if condition {
		return trueVal
	}
	return falseVal
}

//获取提示信息
func getMessage(exception Qexception) string {
	return Ex[exception.DataType] + strconv.Itoa(exception.Status) + ": " + exception.Msg
}

//打印提示信息
func printMsg(exception Qexception) {
	fmt.Println(Ex[exception.DataType] + strconv.Itoa(exception.Status) + ": " + exception.Msg)
}
