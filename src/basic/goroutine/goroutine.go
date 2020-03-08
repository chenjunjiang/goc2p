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

	/*names := []string{"Eric", "Harry", "Robert", "Jim", "Mark"}
	for _, name := range names {
		go func() {
			fmt.Printf("Hello, %s.\n", name)
		}()
	}
	runtime.Gosched()*/
	names := []string{"Eric", "Harry", "Robert", "Jim", "Mark"}
	for _, name := range names {
		fmt.Println(name + "......")
		go func(who string) {
			fmt.Printf("Hello, %s.\n", who)
		}(name)
	}
	runtime.Gosched()
}
