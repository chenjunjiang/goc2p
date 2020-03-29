package main

import (
	"fmt"
	"strconv"
	"time"
)

/**
定时器永远不需要向它的C字段发送第二个元素值。这是因为，定时器一旦到期就意味着在我们重置它之前它无法被再次使用。注意，如果
我们在定时器到期之前停止了它，那么该定时器的字段C也就没有任何缓冲任何元素值了，这之后再去试图从它的C字段中接收元素是不会
有任何结果的。更重要的是，这样做还会使当前Goroutine永远被阻塞。再次强调，重置已被停止的定时器是使它恢复如初的唯一方法。
*/
func main() {
	/*// 初始化一个到期时间距此时的间隔为2秒的定时器
	t := time.NewTimer(2 * time.Second)
	now := time.Now()
	fmt.Printf("Now time: %v.\n", now)
	// 到期事件通过C这个缓冲通道传达的，一旦触及到期时间，定时器就会向自己的C字段发送一个time.Time类型的元素值
	// 这个元素值代表了该定时器的绝对到期时间
	expire := <-t.C
	fmt.Printf("Expiration time: %v.\n", expire)

	//和上面的效果一样
	expire = <-time.After(2 * time.Second)
	fmt.Printf("Expiration time: %v.\n", expire)*/

	/*var t *time.Timer
	f := func() {
		fmt.Printf("Expiration time: %v.\n", time.Now())
		fmt.Printf("C's len: %d\n", len(t.C))
		// 重置之后定时器又恢复执行
		t.Reset(1*time.Second)
	}
	t = time.AfterFunc(1*time.Second, f)
	time.Sleep(200 * time.Second)*/

	/**
	断续器,周期性的传达到期事件,定时器在被重置之前只会传达一次到期事件，而断续器会持续工作直到被停止
	*/
	sign := make(chan byte, 1)
	var ticker *time.Ticker = time.NewTicker(1 * time.Second)
	ticks := ticker.C
	go func() {
		i := 1
		for _ = range ticks {
			fmt.Println("test......" + strconv.Itoa(i))
			if i > 10 {
				ticker.Stop()
				fmt.Println("stop ticker")
				break
			}
			i++
			/**
			close函数是一个内建函数， 用来关闭channel，这个channel要么是双向的， 要么是只写的（chan<- Type）。
			这个方法应该只由发送者调用， 而不是接收者。
			当最后一个发送的值都被接收者从关闭的channel(下简称为c)中接收时,
			接下来所有接收的值都会非阻塞直接成功，返回channel元素的零值。
			如下的代码：
			如果c已经关闭（c中所有值都被接收）， x, ok := <- c， 读取ok将会得到false。
			*/
			// 关闭只读channel会报编译错误，ticks是只读通道
			// close(ticks)
		}
		sign <- 1
	}()
	<-sign

	/**
	在一个定时执行数据修补任务的程序中，为了避免对其它正常的数据库操作产生影响，我们要求两次任务执行之间的最短间隔时间为10分钟，
	可以这样编写代码满足
	*/
	/*var ticker *time.Ticker = time.NewTicker(1 * time.Second)
	ticks := ticker.C
	go func() {
		for _ = range ticks {
			// path代表数据修补任务
			if !path() {
				break
			}
	        // 这条语句的执行会增加下一次迭代的执行延时
			_, ok := <-ticks
			if !ok {
				break
			}
		}
	}()*/
}
