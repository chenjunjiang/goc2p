package main

import (
	"errors"
	"fmt"
)

/**
defer语句被用于预定对一个函数的调用。我们把这类被defer语句调用的函数称为延迟函数。defer语句只能出现在函数或方法的内部。
defer语句的执行时机总是在直接包含它的那个函数（简称外围函数）把流程控制权交给它的调用方的前一刻，无论defer语句出现在外围函数的函数体中的哪一个位置
上。具体分为下面几种情况：
1、当外围函数的函数体中的响应语句全部被正常执行完毕的时候，只有在函数中的所有defer语句都被执行完毕之后该函数才会真正的结束执行。
2、当外围函数的函数体中的return语句被执行的时候，只有在函数中的所有defer语句都被执行完毕之后该函数才会真正地返回。
3、当在外围函数中有运行时恐慌发生的时候，只有在该函数中的所有defer语句都被执行完毕之后该运行时恐慌才会真正地被扩散至该函数的调用方。
这些规则可以和java中的finally比较。
正因为defer语句有这样的特性，所以它成为了执行释放资源或异常处理等收尾任务的首选。
*/
func isPositiveEvenNumber(number int) (result bool) {
	// 多个defer调用满足后进先出，类似栈。
	defer fmt.Println("done.")
	// 调用的函数可以是临时编写的匿名函数
	defer func() {
		fmt.Println("The finishing touches.")
	}()
	if number < 0 {
		panic(errors.New("the number is a negative number")) // 错误字符串不应大写或以标点符号结尾，否则会出现警告
	}
	if number%2 == 0 {
		return true
	}
	return
}

/**
每当defer语句被执行的时候，传递给延迟函数的参数都会以通常的方式被求值。
*/
func begin(funName string) string {
	fmt.Printf("Enter function %s.\n", funName)
	return funName
}

func end(funcName string) string {
	fmt.Printf("Exit function %s.\n", funcName)
	return funcName
}

func record() {
	// begin("record")作为end函数的参数出现
	defer end(begin("record"))
	fmt.Println("In function record.")
}

func printNumbers() {
	for i := 0; i < 5; i++ {
		// 这里会有一个警告,Possible resource leak, 'defer' is called in a for loop.defer的作用域是一个函数，不是一个语句块
		// https://blog.csdn.net/butterfly5211314/article/details/83512711
		defer fmt.Printf("%d ", i)
	}
	fmt.Println("printNumbers......")
}

/**
这个函数被调用之后会输出：5 5 5 5 5
我们说过，在defer语句被执行的时候传递给延迟函数的参数都会被求值。这里的延迟函数是一个没有参数的匿名函数，所以也就没有参数被求值。在for执行完毕的时候，
共有5个相同的延迟函数表达式：
func() {
			fmt.Printf("%d ", i)
		}()
它们在被调用的时候i已经变成5了，所以输出都是5
*/
func printNumbers1() {
	for i := 0; i < 5; i++ {
		defer func() {
			fmt.Printf("%d ", i)
		}()
	}
	fmt.Println("printNumbers1......")
}

func printNumbers2() {
	for i := 0; i < 5; i++ {
		defer func(i int) {
			fmt.Printf("%d ", i)
		}(i)
	}
	fmt.Println("printNumbers2......")
}

func appendNumbers(ints []int) (result []int) {
	result = append(ints, 1)
	defer func() {
		result = append(result, 2)
	}()
	result = append(result, 3)
	defer func() {
		result = append(result, 4)
	}()
	result = append(result, 5)
	defer func() {
		result = append(result, 6)
	}()
	return result
}

/**
如果延迟函数是一个匿名函数，并且在外围函数的声明中存在命名的结果声明，那么在延迟函数中的代码是可以对命名结果的值进行访问和修改的。
*/
func modify(n int) (number int) {
	defer func() {
		number += n
	}()
	number++
	return
}

/**
虽然在延迟函数的声明中可以包含结果声明，但是其返回值会在它被执行完毕时被丢弃。因此，作为惯例，我们在编写延迟函数的声明的时候不会为其添加结果声明。另一方面
，推荐以传参的方式提供延迟函数所需的外部值。
这样调用fmt.Println(modify1(2))之后，输出的值是105，因为外围函数体在执行到return 100的时候，会把number1重新初始化为100（即使前面有number1++）
，但此时还没真正返回，需要等延迟函数执行完毕才会返回，所有最终的number1是105
*/
func modify1(n int) (number1 int) {
	defer func(plus int) (result int) {
		result = n + plus
		number1 += result
		fmt.Println(number1)
		return // 这里的return不会产生任何效果
	}(3)
	number1++
	return 100
}

/**
链式调用，简单来说，在defer x.m1().m2()中，m1会直接被调用，而m2会在最后被调用
下面的例子会输出：
Logger()
do
Log: done

*/
type logger struct {
}

func (l *logger) Print(s interface{}) {
	fmt.Printf("Log: %v\n", s)
}

type customLogger struct {
	l *logger
}

func (f *customLogger) Logger() *logger {
	fmt.Println("Logger()")
	return f.l
}
func do(f *customLogger) {
	// f.Logger()会“先调用”，而f.Print()会“延迟调用”
	defer f.Logger().Print("done")
	fmt.Println("do")
}

func main() {
	fmt.Println(isPositiveEvenNumber(2))
	/**
	Enter function record.
	In function record.
	Exit function record.
	*/
	record()
	printNumbers() // 4 3 2 1 0
	fmt.Println()
	result := appendNumbers([]int{0})
	fmt.Println(result) // [0 1 3 5 6 4 2]
	printNumbers1()     // 5 5 5 5 5
	printNumbers2()     // 4 3 2 1 0
	fmt.Println()
	fmt.Println(modify(2))  // 3
	fmt.Println(modify1(2)) // 105
	var c customLogger = customLogger{l: &logger{}}
	do(&c)
}
