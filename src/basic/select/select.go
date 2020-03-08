package main

import (
	"fmt"
	"time"
)

/**
在运行时系统开始执行select语句的时候，会先对它所有的case中的表达式和通道表达式进行求值。
在执行select语句时，运行时系统会自上而下地判断每个case中的发送或接收操作是否可以立即执行。立即的意思是当前goroutine不会因此操作而被阻塞。当发现第一个
满足条件的case时，运行时系统就会执行该case所包含的语句。其他就会被忽略。如果同时有多个case满足条件，那么运行时系统就会通过一个伪随机的算法决定哪一个
case将会被执行。如果所有case都不满足条件并且没有default case的话，那么当前goroutine就会一直被阻塞于此，直到某一个case中的发送或接收操作可以立即
进行。如果select语句中的所有case右边的通道都是nil，那么当前goroutine就会被永远地阻塞在这个select语句上！所以，通常情况下包含一个default case
总是有必要的。
真正的应用程序中，我们常常需要把select语句放到一个单独的goroutine中，并且为了能连续地接收元素值，应该把select语句包含在一条for循环中。
*/
func main() {
	/*chanCap := 5
	ch7 := make(chan int, chanCap)
	for i := 0; i < chanCap; i++ {
		select {
		case ch7 <- 1:
		case ch7 <- 2:
		case ch7 <- 3:
		}
	}
	for i := 0; i < chanCap; i++ {
		fmt.Printf("%v\n", <-ch7)
	}*/

	sign := make(chan int, 2)
	ch11 := make(chan int, 100)
	go func() {
		var e int
		ok := true
		for {
			select {
			case e, ok = <-ch11:
				if !ok {
					fmt.Println("End.")
					break
				} else {
					fmt.Printf("%d\n", e)
				}
				// 通过调用匿名函数发回一个通道，然后再从该通道中接收元素值赋给ok
			case ok = <-func() chan bool {
				timeout := make(chan bool, 1)
				go func() {
					time.Sleep(time.Millisecond)
					timeout <- false
				}()
				return timeout
			}():
				fmt.Println("Timeout.")
				break
			}
			if !ok {
				break
			}
		}
		sign <- 1
	}()

	go func() {
		for i := 0; i < 50; i++ {
			ch11 <- i
			if i == 31 {
				// 故意模拟超时，上边的代码是大于1毫秒就算超时
				time.Sleep(2 * time.Millisecond)
			}
		}
		close(ch11)
		sign <- 1
	}()

	<-sign
	<-sign
	fmt.Println("all is done.")
}
