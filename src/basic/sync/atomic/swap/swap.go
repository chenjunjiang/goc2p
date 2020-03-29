package main

import (
	"fmt"
	"sync/atomic"
)

/**
与CAS操作不同，原子交换操作不会关心被操作值的旧值。它会直接设置新值。但它又比原子载入操作多做了一步。作为交换，它会返回被操作值的旧值。此类操作比CAS
操作的约束更少，同时又比原子载入操作的功能更强。
*/
func main() {
	var i32 int32 = 1
	oldI32 := atomic.SwapInt32(&i32, 2)
	fmt.Println(oldI32) //1
}
