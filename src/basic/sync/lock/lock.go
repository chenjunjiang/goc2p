package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	/**
	互斥锁
	对于同一个互斥锁的锁定和解锁总是应该成对出现。如果我们锁定了一个已被锁定的互斥锁，那么进行重复锁定操作的Goroutine将会被阻塞，直到该互斥锁回到
	解锁状态。我们一般会在锁定互斥锁之后紧接着就用defer语句来保证该互斥锁及时解锁。虽然互斥锁可以直接在多个Goroutine之间共享，但是还是强烈建议把
	对同一个互斥锁的锁定和解锁操作放在同一个层次的代码块中，例如，在同一个函数或方法中对某个互斥锁进行锁定和解锁。又例如，把互斥锁作为某一结构体类型
	中的字段，以便在该类型的多个方法中使用它。此外，还应该使互斥锁变量的访问权限尽量低。
	*/
	var mutex sync.Mutex
	fmt.Println("Lock the lock.(G0)")
	mutex.Lock()
	fmt.Println("The lock is locked.(G0)")
	for i := 1; i <= 3; i++ {
		go func(i int) {
			fmt.Printf("Lock the lock.(G%d)\n", i)
			mutex.Lock()
			fmt.Printf("The lock is locked.(G%d)\n", i)
		}(i)
	}
	time.Sleep(time.Second)
	fmt.Println("Unlock The lock.(G0)")
	mutex.Unlock()
	fmt.Println("The lock is unlocked.(G0)")
	time.Sleep(time.Second)
}
