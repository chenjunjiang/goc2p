package bmt

import (
	"fmt"
	"strconv"
	"testing"
)

/**
基准测试，是一种测试代码性能的方法，比如你有多种不同的方案，都可以解决问题，那么到底是那种方案性能更好呢？这时候基准测试就派上用场了。
基准测试主要是通过测试CPU和内存的效率问题，来评估被测试代码的性能，进而找到更好的解决方案。比如链接池的数量不是越多越好，那么哪个值才是最优值呢，
这就需要配合基准测试不断调优了。
基准测试有以下规则：
    基准测试的代码文件必须以_test.go结尾
    基准测试的函数必须以Benchmark开头，必须是可导出的
    基准测试函数必须接受一个指向Benchmark类型的指针作为唯一参数
    基准测试函数不能有返回值
    b.ResetTimer是重置计时器，这样可以避免for循环之前的初始化代码的干扰
    最后的for循环很重要，被测试的代码要放到循环里
    b.N是基准测试框架提供的，表示循环的次数，因为需要反复调用测试的代码，才可以评估性能

运行基准测试也要使用go test命令，不过我们要加上-bench=标记，它接受一个表达式作为参数，匹配基准测试的函数，.表示运行所有基准测试。
chenjunjiang@chenjunjiang-B85-HD3:~/go_workspace/goc2p/src/basic/testing/bmt$ go test . -bench=.
goos: linux
goarch: amd64
pkg: basic/testing/bmt
BenchmarkSprintf-4   	11832226	        99.2 ns/op
PASS
ok  	basic/testing/bmt	1.279s
chenjunjiang@chenjunjiang-B85-HD3:~/go_workspace/goc2p/src/basic/testing/bmt$ go test . -bench=Sprintf
goos: linux
goarch: amd64
pkg: basic/testing/bmt
BenchmarkSprintf-4   	11646510	        99.4 ns/op
PASS
ok  	basic/testing/bmt	1.264s

因为默认情况下go test 会运行单元测试，为了防止单元测试的输出影响我们查看基准测试的结果，可以使用-run=匹配一个从来没有的单元测试方法，
过滤掉单元测试的输出，我们这里使用none，因为我们基本上不会创建这个名字的单元测试方法。
chenjunjiang@chenjunjiang-B85-HD3:~/go_workspace/goc2p/src/basic/testing/bmt$ go test . -bench=Sprintf -run=none
goos: linux
goarch: amd64
pkg: basic/testing/bmt
BenchmarkSprintf-4   	12022448	        99.0 ns/op
PASS
ok  	basic/testing/bmt	1.293s

下面着重解释下说出的结果，看到函数后面的-4了吗？这个表示运行时对应的GOMAXPROCS的值。接着的12022448表示运行for循环的次数，
也就是调用被测试代码的次数，最后的99.0 ns/op表示每次需要花费99纳秒。
以上是测试时间默认是1秒，也就是1秒的时间，调用12022448次，每次调用花费99纳秒。如果想让测试运行的时间更长，可以通过-benchtime指定，比如3秒。
chenjunjiang@chenjunjiang-B85-HD3:~/go_workspace/goc2p/src/basic/testing/bmt$ go test . -bench=Sprintf -benchtime=3s -run=none
goos: linux
goarch: amd64
pkg: basic/testing/bmt
BenchmarkSprintf-4   	33834154	        99.7 ns/op
PASS
ok  	basic/testing/bmt	3.484s
可以发现，我们加长了测试时间，测试的次数变多了，但是最终的性能结果：每次执行的时间，并没有太大变化。一般来说这个值最好不要超过3秒，意义不大。

这个例子其实是一个int类型转为string类型的例子，标准库里还有几种方法，我们看下哪种性能更加。
*/
func BenchmarkSprintf(b *testing.B) {
	num := 10
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("%d", num)
	}
}

func BenchmarkFormat(b *testing.B) {
	num := int64(10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		strconv.FormatInt(num, 10)
	}
}

func BenchmarkItoa(b *testing.B) {
	num := 10
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		strconv.Itoa(num)
	}
}

/**
chenjunjiang@chenjunjiang-B85-HD3:~/go_workspace/goc2p/src/basic/testing/bmt$ go test -bench=. -run=none
goos: linux
goarch: amd64
pkg: basic/testing/bmt
BenchmarkSprintf-4   	11760518	        99.2 ns/op
BenchmarkFormat-4    	336680968	         3.50 ns/op
BenchmarkItoa-4      	342064897	         3.52 ns/op
PASS
ok  	basic/testing/bmt	4.372s
 从结果上看strconv.FormatInt函数是最快的，其次是strconv.Itoa，然后是fmt.Sprintf最慢。那么最后一个为什么这么慢的，
我们再通过-benchmem找到根本原因。
chenjunjiang@chenjunjiang-B85-HD3:~/go_workspace/goc2p/src/basic/testing/bmt$ go test -bench=. -benchmem -run=none
goos: linux
goarch: amd64
pkg: basic/testing/bmt
BenchmarkSprintf-4   	11742652	        98.9 ns/op	      16 B/op	       2 allocs/op
BenchmarkFormat-4    	338487405	         3.52 ns/op	       0 B/op	       0 allocs/op
BenchmarkItoa-4      	340388349	         3.53 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	basic/testing/bmt	4.379s
-benchmem可以提供每次操作分配内存的次数，以及每次操作分配的字节数。allocs/op表示每次操作从堆上分配内存的次数。B/op表示每次操作分配的字节数。
从结果我们可以看到，性能高的两个函数，每次操作都是进行0次内存分配，而最慢的那个要分配2次；性能高的每次操作分配0个字节内存，
而慢的那个函数每次需要分配16字节的内存。从这个数据我们就知道它为什么这么慢了，内存分配和占用都太高。
在代码开发中，对于我们要求性能的地方，编写基准测试非常重要，这有助于我们开发出性能更好的代码。不过性能、可用性、复用性等也要有一个相对的取舍，
不能为了追求性能而过度优化。
*/
