package main

import (
	"fmt"
	"time"
)

/**
在接收元素值的时候，如果当前通道中没有任何元素值，当前Goroutine会被阻塞于此。如果在进行接收操作之前或过程当中该通道被关闭了，那么该操作会立即被结束。
一个通道的缓冲容量是固定的。因此，在通道里面的元素数量达到最大后，当某个Goroutine再向还通道发送值时，该Goroutine会被阻塞，直到该通道中有足够的空间
容纳该元素为止。
通过 len 函数可以获得 chan 中的元素个数，通过 cap 函数可以得到 channel 的缓存长度。
*/
type Person struct {
	Name    string
	Age     uint8
	Address Addr
}

type Addr struct {
	city     string
	district string
}

type PersonHandler interface {
	Batch(origs <-chan Person) <-chan Person
	Handle(orig *Person)
}

type PersonHandlerImpl struct {
}

func (handler PersonHandlerImpl) Batch(origs <-chan Person) <-chan Person {
	dests := make(chan Person, 100)
	go func() {
		/**
		当通道中没有任何元素的时候，当前goroutine会阻塞在range处。在相关的通道被关闭后，若通道中已无元素值或当前goroutine正阻塞于此，则这条for
		语句的执行会 立即结束;而当此时的通道中还有遗留元素值时，运行时系统会等for语句把他它们全部接收后再结束该语句的执行。
		*/
		for p := range origs {
			handler.Handle(&p)
			dests <- p
		}
		fmt.Println("All the information has been handled.")
		close(dests)
	}()
	return dests
}

func (handler PersonHandlerImpl) Handle(orig *Person) {
	if orig.Address.district == "wuhou" {
		orig.Address.district = "jinjiang"
	}
}

var personTotal = 200
var persons = make([]Person, personTotal)
var personCount int

/**
初始化persons
*/
func init() {
	for i := 0; i < 200; i++ {
		name := fmt.Sprintf("%s%d", "P", i)
		p := Person{name, 32, Addr{"chengdu", "wuhou"}}
		persons[i] = p
	}
}

func getPersonHandler() PersonHandler {
	return PersonHandlerImpl{}
}

func savePerson(dest <-chan Person) <-chan byte {
	sign := make(chan byte, 1)
	go func() {
		for {
			p, ok := <-dest
			if !ok {
				fmt.Println("All the information has been saved.")
				sign <- 0
				break
			}
			savePerson1(p)
		}
	}()
	return sign
}

func fetchPerson(origs chan<- Person) {
	origsCap := cap(origs)
	// 通过容量判断当前通道的类型是缓冲通道还是非缓冲通道
	buffered := origsCap > 0
	goTicketTotal := origsCap / 2
	fmt.Println("goTicketTotal=", goTicketTotal)
	goTicket := initGoTicket(goTicketTotal)
	fmt.Println("buffered start=", len(goTicket))
	go func() {
		for {
			p, ok := fetchPerson1()
			// 下面注释处这种写法会导致origs不能被关闭，那么在其它地方使用origs和dests通道的地方就会导致死锁
			/*if !ok {
				if !buffered || len(goTicket) == goTicketTotal {
					break
				}
				time.Sleep(time.Nanosecond)
				fmt.Println("All the information has been fetched.")
				close(origs)
				break
			}
			origs <- p*/
			fmt.Println("ok=", ok)
			if !ok {
				for {
					fmt.Println("bufferedxxxx=", len(goTicket))
					if !buffered || len(goTicket) == goTicketTotal {
						fmt.Println("buffered=", len(goTicket))
						break
					}
					time.Sleep(time.Nanosecond)
				}
				fmt.Println("All the information has been fetched.")
				close(origs)
				break
			}
			/**
			通过goTicket通道来限制启动Goroutine的数量，只有接收元素值<-goTicket不被阻塞了，才能继续下面的新建Goroutine的动作
			*/
			if buffered {
				// 每次循环都会从goTicket接收元素值
				<-goTicket
				go func() {
					origs <- p
					// 每个goroutine执行完之前又向goTicket发送元素值，只有所有的goroutine执行完毕才能保证上面的len(goTicket) == goTicketTotal 条件成立
					goTicket <- 1
				}()
			} else { // 如果origs是非缓冲通道就没必要并发地发送人员信息了，因为非缓冲通道只能同步地传递元素值。在接收完成之前，发送操作是无法完成的
				origs <- p
			}
		}
	}()
}

func fetchPerson1() (Person, bool) {
	if personCount < personTotal {
		p := persons[personCount]
		personCount++
		return p, true
	}
	return Person{}, false
}

func initGoTicket(total int) chan byte {
	var goTicket chan byte
	if total == 0 {
		return goTicket
	}
	goTicket = make(chan byte, total)
	for i := 0; i < total; i++ {
		goTicket <- 1
	}
	return goTicket
}

func savePerson1(p Person) bool {
	return true
}

func main() {
	/**
	在发送过程中进行的元素值复制属于完全复制，通道中的元素值不会受外界影响
	*/
	/*var personChan = make(chan Person, 1)
	p1 := Person{"chenjj", 31, Addr{"chengdu", "wuhou"}}
	fmt.Printf("p1 (1): %v\n", p1)
	personChan <- p1
	p1.Address.district = "jinjiang"
	fmt.Printf("p1 (2): %v\n", p1)
	p1Copy := <-personChan
	fmt.Printf("p1_copy: %v\n", p1Copy)*/

	/**
	无论怎样都不应该在接收端关闭通道。因为那里我们无法判断发送端是否还会向通道发送元素值。然而，我们在发送端调用close关闭通道却不会对接收端接收元素造成
	任何影响，接收端还是会把通道中的所有元素值都接收到。这样保证了通道的安全性。
	对于同一通道，只允许关闭一次。
	*/
	/*ch := make(chan int, 5)
	sign := make(chan byte, 2)
	go func() {
		for i := 0; i < 5; i++ {
			ch <- i
			time.Sleep(1 * time.Second)
		}
		close(ch)
		fmt.Println("The channel is closed.")
		sign <- 0
	}()
	go func() {
		for {
			e, ok := <-ch
			fmt.Printf("%d (%v)\n", e, ok)
			// 通道已被关闭
			if !ok {
				break
			}
			time.Sleep(2 * time.Second)
		}
		fmt.Println("Done.")
		sign <- 1
	}()
	// 当sign通道中没有元素的话，接收操作就会阻塞当前Goroutine
	<-sign
	<-sign*/

	/**
	通道所允许的数据传递方向是它类型的一部分。对于双通道类型而言，方向的不同就意味着它们类型的不同。也就是说，元素类型相同的双向通道、发送通道和接收
	通道都属于不同的类型。因此，我们只能利用函数声明来约束通道的方向。比如，利用函数的参数声明把函数调用方所持有的双向通道转换为单向通道，并提供给函数
	内部使用。又比如，利用函数的结果声明把函数内部所持有的双向通道转换为单向通道，并提供给函数的调用方使用。注意，即使利用函数声明转换通道类型，也无法
	把单向通道转换为双向通道。并且也不能改变单向通道的方向。
	*/
	handler := getPersonHandler()
	origs := make(chan Person, 100)
	dets := handler.Batch(origs)
	fetchPerson(origs)
	sign := savePerson(dets)
	<-sign
}
