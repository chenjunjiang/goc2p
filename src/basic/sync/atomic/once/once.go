package main

import (
	"fmt"
	"time"
)

/**
我们对sync.Once类型值的指针方法Do的有效调用次数永远是1。也就是说 ，无论我们调用这个方法多少次，都只有第一次调用是有效的。
典型的应用场景就是仅需要执行一次的任务。比如，数据库连接池的初始化任务，一些需要持续运行的实时检测任务等等。
*/
func main() {
	onceDo()
	/**
	Received a signal.
	Timeout!
	Timeout!
	Num: 2.

	*/
}

func onceDo() {
	var num int
	sign := make(chan bool)
	//var once sync.Once
	f := func(ii int) func() {
		return func() {
			num = num + ii*2
			sign <- true
		}
	}
	for i := 0; i < 3; i++ {
		fi := f(i + 1)
		//go once.Do(fi)
		go fi()
	}
	for j := 0; j < 3; j++ {
		select {
		case <-sign:
			fmt.Println("Received a signal.")
		case <-time.After(100 * time.Millisecond):
			fmt.Println("Timeout!")
		}
	}
	fmt.Printf("Num: %d.\n", num)
}
