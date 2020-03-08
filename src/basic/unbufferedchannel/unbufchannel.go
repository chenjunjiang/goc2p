package main

import (
	"fmt"
	"time"
)

/**
非缓冲通道有三个特别之处：
1、向此类通道发送元素值的操作会被阻塞，直到至少有一个针对该通道的接收操作开始进行为止。
2、从此类通道接收元素值的操作会被阻塞，直到至少有一个针对该通道的发送操作开始进行为止。
3、针对非缓冲通道的接收操作会在与之对应的发送操作完成之前完成。
在缓冲通道中，由于元素值的传递是异步的，所以发送操作在成功向通道发送元素值之后就会立即结束。然而，针对非缓冲通道的操作在这方面的表现正好相反。发送操作
向非缓冲通道发送元素值的时候，会等待能够接收该元素值的那个接收操作。并且，只有确保该元素值被成功接收，它才会真正的完成执行。
*/
func main() {
	// 可以利用非缓冲通道的特性来实现多个Goroutine之间的同步
	unbufChan := make(chan int)
	go func() {
		fmt.Println("Sleep a second...")
		time.Sleep(time.Second)
		// 接收操作完成后下面的发送操作才会执行完成
		num := <-unbufChan
		fmt.Printf("Received a integer %d.\n", num)
	}()
	num := 1
	fmt.Printf("Send integer %d...\n", num)
	unbufChan <- num
	fmt.Println("Done.")
}
