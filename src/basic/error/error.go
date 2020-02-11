package main

import (
	"errors"
	"fmt"
)

/**
为了使编程人员能够在自己的程序中报告运行期间的、不可恢复的错误状态，Go语言使用内建函数panic来停止当前的控制流程的执行并报告一个运行时恐慌。
*/

func outerFunc() {
	innerFunc()
}

func innerFunc() {
	panic(errors.New("A intended fatal error"))
}

func fetchDemo() {
	defer func() {
		if v := recover(); v != nil {
			fmt.Printf("Recovered a panic. [index=%d]\n", v)
		}
	}()
	ss := []string{"A", "B", "C"}
	fmt.Printf("Fetch the elements in %v one by one...\n", ss)
	r := fetchElement(ss, 0)
	fmt.Println(r)
	fmt.Println("The elements fetching is done.")
}

func fetchElement(ss []string, index int) (element string) {
	if index >= len(ss) {
		fmt.Printf("Occur a panic! [index=%d]\n", index)
		panic(index)
	}
	fmt.Printf("Fetching the element... [index=%d]\n", index)
	element = ss[index]
	defer fmt.Printf("The element is \"%s\". [index=%d]\n", element, index)
	fetchElement(ss, index+1)
	/*if index < 2 {
		fetchElement(ss, index+1)
	}*/
	return
}

func main() {
	// outerFunc()
	//myIndex := 4
	//ia := [3]int{1, 2, 3}
	/**
	这句代码会引发一个运行时恐慌，因为它造成了数组越界，这个运行时恐慌是由运行时系统报告的。它相当于我们显示地调用panic函数并传入一个 runtime.Error
	类型的参数值。
	*/
	//_ = ia[myIndex]
	fetchDemo()
}
