package main

import (
	"fmt"
	"sync/atomic"
)

/**
Compare And Swap
*/
func main() {
	var i32 int32 = 3
	atomic.CompareAndSwapInt32(&i32, 3, 6)
	fmt.Println(i32) // 6
	// 在多Goroutine环境下，一般都是利用for循环进行多次尝试，比如使用cas来实现加法
	addValue(&i32, 4)
	fmt.Println(i32)
}

// 这里的value必须是指针类型的，否则改变的只是复制的值，调用方的值不会变
func addValue(value *int32, delta int32) {
	for {
		v := *value
		if atomic.CompareAndSwapInt32(value, v, v+delta) {
			break
		}
	}
}
