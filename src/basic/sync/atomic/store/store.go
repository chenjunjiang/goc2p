package main

import (
	"fmt"
	"sync/atomic"
)

/**
在原子地存储某个值的过程中，任何CPU都不会进行针对同一个值的读或写操作。就不会出现针对此值的读操作被并发进行而读到修改了一半的值的情况(32位操作系统
上操作64位数字的时候可能会出现这种情况)。
原子的值存储操作总会成功，因为它并不关心被操作值的旧值是什么。
*/
func main() {
	var i32 int32 = 1
	atomic.StoreInt32(&i32, 2)
	fmt.Println(i32) // 2
}
