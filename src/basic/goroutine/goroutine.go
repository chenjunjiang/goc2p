package main

import (
	"fmt"
	"runtime"
)

/**
封装main函数的Goroutine是Go语言运行时系统创建的第一个Goroutine（也被称为主Goroutine）。主Goroutine是在runtime.m0上被运行的。实际上，在
runtime.m0运行完runtime.g0中的引导程序之后，会接着运行主Goroutine。
*/
func main() {
	/*name:="Eric"
	go func() {
		fmt.Printf("Hello, %s.\n", name) // Hello, Harry.
	}()
	name="Harry"
	// 让其它goroutine有机会运行，但不保证其它goroutine一定比当前goroutine先执行
	runtime.Gosched()*/

	names := []string{"Eric", "Harry", "Robert", "Jim", "Mark"}
	for _, name := range names {
		// 这里如果不暂停，下边的Goroutine可能全部输出的都是Hello, Mark.不要对Goroutine的执行时机做任何假设。什么时候执行是调度器决定的。
		// time.Sleep(100 * time.Millisecond)
		go func() {
			fmt.Printf("Hello, %s.\n", name)
		}()
	}
	runtime.Gosched()
	/*names := []string{"Eric", "Harry", "Robert", "Jim", "Mark"}
	for _, name := range names {
		fmt.Println(name + "......")
		go func(who string) {
			fmt.Printf("Hello, %s.\n", who)
		}(name)
	}
	runtime.Gosched()*/
}
