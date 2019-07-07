package main

import (
	"fmt"
	"os"
	"strconv"
)

/**
程序的命令行参数可从os包的Args变量获取
对string类型，+运算符连接字符串
chenjunjiang@chenjunjiang-B85-HD3:~/go_workspace/goc2p/src/basic/args$ go run main.go 1 3 - x ?
chenjunjiang@chenjunjiang-B85-HD3:~/go_workspace/goc2p/src/basic/args$ go build
chenjunjiang@chenjunjiang-B85-HD3:~/go_workspace/goc2p/src/basic/args$ ./args 1 2
1 2
./args
 */
func main() {
	/*var s, sep string
	for i := 1; i < len(os.Args); i++ {
		s += sep + os.Args[i]
		sep = " "
	}
	fmt.Println(s)
	// 被执行命令本身的名字
	fmt.Println(os.Args[0])*/

	// range迭代的时候会返回两个值分别是数据的索引和值
	/**
	这个例子不需要索引，但range的语法要求, 要处理元素, 必须处理索引。一种思路是把索引赋值给一个临时变量, 如temp, 然后忽略它的值，
	但Go语言不允许使用无用的局部变量（local variables），因为这会导致编译错误。Go语言中这种情况的解决方法是用空标识符（blank identifier）
	，即_（也就是下划线）。空标识符可用于任何语法需要变量名但程序逻辑不需要的时候
	 */
	/*s, sep := "", ""
	for _, arg := range os.Args[1:] {
		s += sep + arg
		sep = " "
	}
	fmt.Println(s)*/

	/**
	上面的两种方式每次循环迭代字符串s的内容都会更新。+=连接原字符串、空格和下个参数，产生新字符串, 并把它赋值给s。s原来的内容已经不再使用，
	将在适当时机对它进行垃圾回收。如果连接涉及的数据量很大，这种方式代价高昂。一种简单且高效的解决方案是使用strings包的Join函数
	 */
	// fmt.Println(strings.Join(os.Args[1:], " "))

	// 打印每个参数的索引和值，每个一行
	for index, arg := range os.Args[1:] {
		s, sep := "", " "
		s += strconv.Itoa(index) + sep + arg
		fmt.Println(s)
		string := strconv.FormatInt(int64(index), 10)
		// strconv.Atoi(string)
		fmt.Println(string)
	}
}
