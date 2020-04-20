package main

import (
	"fmt"
	"sync/atomic"
)

/**
被用于进行增或减的原子操作的函数名称都以"Add"为前缀，并后跟针对具体类型的名称。
*/
func main() {
	var i32 int32 = 1
	// 第一个参数之所以必须是指针类型的值，是因为该函数需要获得被操作值在内存中的存放位置，以便施加特殊的CPU指令。
	newI32 := atomic.AddInt32(&i32, 3)
	fmt.Println(newI32) // 4
	fmt.Println(i32)    // 4
	var i64 int64 = 5
	atomic.AddInt64(&i64, -3)
	fmt.Println(i64) // 2
	var uI64 uint64 = 5
	// atomic.AddUint64(&uI64, -3) // constant -3 overflows uint64
	var NN int64 = -3
	// 针对无符号数的减法，需要先定义一个变量NN再操作，利用了二进制补码的特性
	atomic.AddUint64(&uI64, ^uint64(-NN-1)) // 把uI64减3
	fmt.Println(uI64)                       // 2
}
