package util

import (
	"fmt"
	"runtime"
)

//GetCallStackInfo 返回函数调用堆栈信息
func GetCallStackInfo(callStack int) (fileName string, codeLine int, functionName string) {
	pc, file, line, ok := runtime.Caller(callStack)
	if ok {
		fileName = file
		codeLine = line
		fun := runtime.FuncForPC(pc)
		if nil != fun {
			functionName = fun.Name()
		}
	}
	return
}

//GetCurrentCallStackInfo 返回函数调用堆栈信息
func GetCurrentCallStackInfo() (fileName string, codeLine int, functionName string) {
	return GetCallStackInfo(2) //需要返回 调用GetCurrentCallStackInfo的位置的信息
}

// CatchError 捕获异常
func CatchError() {
	if err := recover(); err != nil {
		fmt.Printf("recover error:%v\n", err)
	}
}
