package et

import "fmt"

/**
样本测试函数的名称需要以"Example"作为开始。并且，在这类函数的函数体的最后还可以有若干个注释行。这些注释行的作用是，比较在该测试函数被执行期间，标准
输出上出现的内容是否与预期的相符。要想使这些注释行能被正确地解析，需要满足以下几个条件：
1、这些注释行必须出现在函数体的末尾，且在它们和作为当前函数体结束符的"}"之前没有任何代码。
2、在第一行注释中，紧跟在单行注释前导符//之后永远应该是Output:。
3、在Output:右边的内容以及后续注释行中的内容都分别代表了标准输出中的一行内容。
样本函数的运行：
chenjunjiang@chenjunjiang-B85-HD3:~/go_workspace/goc2p/src/basic/testing/et$ go test -v
=== RUN   ExampleHello
--- PASS: ExampleHello (0.00s)
PASS
ok  	basic/testing/et	0.001s

chenjunjiang@chenjunjiang-B85-HD3:~/go_workspace/goc2p/src/basic/testing/et$ go test -v
=== RUN   ExampleHello
--- FAIL: ExampleHello (0.00s)
got:
Hello,Golang~
want:
Hello,Golang1~
FAIL
exit status 1
FAIL	basic/testing/et	0.001s

如果一个测试函数中没有任何样本注释行，那么这个函数仅仅会被编译而不会被执行。
样本测试函数的命名规则：
1、当被测试对象是整个代码包时，名称应该是Example。
2、当被测试对象是一个函数时，对于函数F，样本测试函数的名称应该是ExampleF。
3、当被测试对象是某个类型中的一个方法时，对于类型T，样本测试函数的名称应该是ExampleT。
3、当被测试对象是某个类型中的一个方法时，对于类型T的方法M，样本测试函数的名称应该是ExampleT_M。
*/
/*func ExampleHello() {
	fmt.Println("Hello,Golang~")
	// Output: Hello,Golang1~
}*/
func ExampleHello() {
	for i := 0; i < 3; i++ {
		fmt.Println("Hello,Golang~")
	}
	// Output: Hello,Golang~
	// Hello,Golang~
	// Hello,Golang~
}
