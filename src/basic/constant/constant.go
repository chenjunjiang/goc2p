package main

import "fmt"

/**
在Go语言中，常量总会在编译期间被创建，即使它们作为局部变量被定义在了函数内部。也正因为如此，常量只能由字面常量或常量表达式来赋值。
Go语言的常量可以被划分为布尔常量、rune常量（字符常量）、整数常量、浮点数常量、复数常量和字符串常量。
*/
const untypedConstant = 10.0     // 无类型常量
const typedConstant int64 = 1024 // 有类型常量

// 把多个常量的声明拆分多行
const (
	utc1 = 6.3
	utc2 = false
	utc3 = "C"
)

/**
在带圆括号的常量声明语句中，有时候我们并不需要显示地对所有的常量进行赋值。被省略了赋值的常量实际上还是有值的，只不过这个值是被隐含地赋予的。它们的值及其
类型都会与在其上面、最近的且被显示赋值的那个常量相同。
*/
const (
	utc4 = 6.3
	utc5 = "D"
	utc6
	utc7
)

const (
	utc8        = 6.3
	utc9, utc10 = false, "E"
	utc11, utc12
)

/**
在常量声明语句中，iota代表了连续的、无类型的整数常量。它第一次出现在一个以const开始的常量声明语句中的时候总会表示整数常量0。随着在同一条常量声明语句
中包含iota的常量声明的数量的增加，iota所表示的整数值也会递增。
*/
// x和y的值都是整数常量0，因为它们出现在了两条常量声明语句中。它们之间并不存在递增关系。
const x = iota
const y = iota
const (
	a = iota
	b
	c
)
const (
	u = 1 << iota // 1<<0
	v             // 1<<1
	w             // 1<<2
)

// 可以利用空标识符"_"来跳过iota表示的递增序列中的某个或某些值
const (
	e, f = iota, 1 << iota // 0,1<<0
	_, _
	g, h // 2, 1<<2
	i, j // 3,1<<3
)

func main() {
	fmt.Println(untypedConstant)
	fmt.Println(typedConstant)
	fmt.Println(utc3)
	fmt.Println(utc6)  // D
	fmt.Println(utc7)  // D
	fmt.Println(utc11) // false
	fmt.Println(utc12) // E
	fmt.Println(a)     // 0
	fmt.Println(b)     // 1
	fmt.Println(c)     // 2
	s := "中国"
	fmt.Println(string(s[1]))         // 乱码, s[1]取的是第2个字节，string底层是byte数组
	fmt.Println(string([]rune(s)[1])) // 国

}
