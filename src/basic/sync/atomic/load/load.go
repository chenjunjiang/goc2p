package main

import (
	"fmt"
	"sync/atomic"
)

/**
在32位计算架构的计算机上写入一个64位的整数，如果这个写操作未完成的时候，有一个读操作被并发地执行了，那么这个读操作很可能会读取到一个只被修改了
一半的数据。为了原子地读取某个值，可以使用sync/atomic代码包提供的以"Load"为前缀的函数。
*/
func main() {
	var i32 int32 = 1
	addValue(&i32, 1)
	fmt.Println(i32) // 2
}

func addValue(value *int32, delta int32) {
	for {
		/**
		使用atomic.LoadInt32的含义是原子地读取变量value的值并把它赋给变量v。原子读取就意味着在读取value的同时，当前计算机中的任何CPU都不会
		进行其它针对此值的读或写操作。这样的约束是受到底层硬件支持的。
		注意：虽然我们在这里原子地载入value的值，但是后面的CAS操作仍然是必要的。因为把值赋给v和if语句并不会原子执行。在它们被执行期间，CPU仍然
		可能进行其它针对value的读或写操作。
		*/
		v := atomic.LoadInt32(value)
		if atomic.CompareAndSwapInt32(value, v, v+delta) {
			break
		}
	}
}
