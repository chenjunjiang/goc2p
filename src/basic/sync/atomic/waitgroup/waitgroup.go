package main

import (
	"fmt"
	"sync"
	"time"
)

/**
假设我们的程序启动了4个Goroutine，分别是G1、G2、G3和G4。其中，G2、G3和G4是由G1中的代码启用并被用于执行某些特定的任务的。G1在启用这3个Goroutine
之后要等待这些特定的任务完成。我们有两种方案：通道和WaitGroup。
通道的实现中包含了很多专门为并发安全地传递数据而建立的数据结构和算法，原则上，我们不应该把通道当做互斥锁或信号灯来使用。
*/
func main() {
	/*sign := make(chan byte, 3)
	go func() {
		fmt.Println("执行G2任务")
		time.Sleep(time.Second)
		sign <- 2
	}()
	go func() {
		fmt.Println("执行G3任务")
		time.Sleep(time.Second)
		sign <- 3
	}()
	go func() {
		fmt.Println("执行G4任务")
		time.Sleep(time.Second)
		sign <- 4
	}()
	for i := 0; i < 3; i++ {
		fmt.Printf("G%d is ended.\n", <-sign)
	}*/
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		fmt.Println("执行G2任务")
		time.Sleep(time.Second)
		wg.Done()
	}()
	go func() {
		fmt.Println("执行G3任务")
		time.Sleep(time.Second)
		wg.Done()
	}()
	go func() {
		fmt.Println("执行G4任务")
		time.Sleep(time.Second)
		wg.Done()
	}()
	wg.Wait()
	fmt.Println("G2, G3 and G4 are ended.")
}
