package main

import (
	"fmt"
	"runtime"
)

func init()  {// 包初始化函数
	fmt.Printf("Map: %v\n", m) // 先格式化再打印
	info = fmt.Sprintf("OS: %s, Arch: %s", runtime.GOOS, runtime.GOARCH)
}

var m map[int]string = map[int]string{1:"A",2:"B",3:"C"}
var info string

func main() {
	
}
