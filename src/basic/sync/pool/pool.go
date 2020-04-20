package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
	"sync/atomic"
)

/**
我们通常用golang来构建高并发场景下的应用，但是由于golang内建的GC机制会影响应用的性能，为了减少GC，golang提供了对象重用的机制，
也就是sync.Pool对象池(临时对象池)。 sync.Pool是可伸缩的，并发安全的。其大小仅受限于内存的大小，可以被看作是一个存放可重用对象的值的容器。
设计的目的是存放已经分配的但是暂时不用的对象，在需要用到的时候直接从pool中取。
临时对象池会专门为每一个与操作它的Goroutine相关联的P都生成一个本地池。

Put方法逻辑：
1、如果放入的值为空，直接return。
2、检查当前goroutine的是否设置对象池私有值，如果没有则将x赋值给其私有成员，并将x设置为nil。
3、如果当前goroutine私有值已经被设置，那么将该值追加到共享列表。
Get方法逻辑：
1、尝试从本地P对应的那个本地池中获取一个对象值, 并从本地池冲删除该值。
2、如果获取失败，那么从共享池中获取, 并从共享队列中删除该值。
3、如果获取失败，那么从其他P的共享池中偷一个过来，并删除共享池中的该值(p.getSlow())。
4、如果仍然失败，那么直接通过New()分配一个返回值，注意这个分配的值不会被放入池中。New()返回用户注册的New函数的值，如果用户未注册New，那么返回nil。

临时对象池对垃圾回收友好，垃圾回收的执行一般会使临时对象池中的对象值被全部移除，也就是说，即使我们永远不会显示地从临时对象池取走某一个对象值，该对象值
也不会永远待在临时对象池中。它的生命周期取决于垃圾回收任务下一次的执行是时间。
临时对象池中的任何对象值都有可能在任何时候被移除掉，并且根本不会通知该池的调用方。这种情况常常会发生在垃圾回收器即将开始回收内存垃圾的时候。如果这时临时
对象池中的某个对象值仅被该池引用，那么它可能会在垃圾回收的时候被回收掉。
根据上面的说法，Golang的对象池严格意义上来说是一个临时的对象池，适用于储存一些会在goroutine间分享的临时对象。主要作用是减少GC，提高性能。
在Golang中最常见的使用场景是fmt包中的输出缓冲区。
*/
func main() {
	// 禁用GC，并保证在main函数执行结束之前恢复GC
	// defer debug.SetGCPercent(debug.SetGCPercent(-1))
	var count int32
	newFunc := func() interface{} {
		return atomic.AddInt32(&count, 1)
	}
	// 通过New去定义你这个池子里面放的究竟是什么东西，在这个池子里面你只能放一种类型的东西。
	pool := sync.Pool{New: newFunc}
	v1 := pool.Get()
	fmt.Printf("v1: %v\n", v1)
	pool.Put(newFunc())
	fmt.Printf("count: %v\n", count)
	v2 := pool.Get()
	fmt.Printf("v2: %v\n", v2)
	// 垃圾回收对临时对象池的影响
	debug.SetGCPercent(100)
	runtime.GC()
	v3 := pool.Get()
	fmt.Printf("v3: %v\n", v3)
	pool.New = nil
	v4 := pool.Get()
	fmt.Printf("v4: %v\n", v4)
}
