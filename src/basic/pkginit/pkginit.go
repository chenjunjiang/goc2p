package main

import (
	"fmt"
	"runtime"
)

/**
  同一个代码包中可以存在多个代码包初始化函数，甚至代码包内的每一个源码文件都可以定义多个代码包初始化函数。Go不会保证同一个代码包中多个代码包初始化函数
的执行顺序。此外，被导入的代码包的初始化函数总是先会执行。
  init()在main函数执行之前执行，并且只会执行一次。
*/
func init() { // 包初始化函数
	fmt.Printf("Map: %v\n", m) // 先格式化再打印
	info = fmt.Sprintf("OS: %s, Arch: %s", runtime.GOOS, runtime.GOARCH)
}

// 当前代码包中的所有全局变量的初始化会在代码包初始化函数执行前完成。这就避免了在代码包初始化函数对某个变量进行赋值之后又被该变量声明中赋予的值覆盖掉的问题
var m map[int]string = map[int]string{1: "A", 2: "B", 3: "C"}

// info没有显示赋值，它被赋予string类型的零值""(空字符串)
var info string

func main() { // 命令源码文件必须有的入口函数
	fmt.Println(info) // 打印变量info
}
