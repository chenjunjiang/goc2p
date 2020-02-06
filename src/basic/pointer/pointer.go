package main

import (
	"fmt"
	"unsafe"
)

/**
指针是一个代表着某个内存地址的值。这个内存地址往往是在内存中存储的另一个变量的值的起始位置。
Go语言中有一个专门用于存储内存地址的类型uintptr，和int和uint类型一样，属于数值类型。它的值是一个能够保存一个指针值的32位或64位无符号整数。
也可以说，它的值是指针类型值的位模式形式。
1、类型表示法
在任何一个有效的数据类型的左边插入符号*来得到与之对应的指针类型。比如：一个元素类型为int的切片类型对应的指针类型是*[]int。
2、值表示法
如果一个变量v的值是可寻址的，那么我们可以使用取址操作符&取出与这个值对应的指针值。表达式&v就代表了指向变量v的值的指针值。
如果某个值确实被存储在了计算机内存中，并且有一个内存地址可以代表这个值在内存中存储的起始位置，那么我们就可以说这个值以及代表它的变量是可寻址的。

指针类型理所当然属于引用类型 ，它的零值是nil。
在代码包unsafe中， 有一个名为ArbitraryType的类型。它是int类型的一个别名类型。但是，它实际上可以代表任意的Go语言表达式的结果类型。unsafe包还声明
了一个名为Pointer的类型。unsafe.Pointer类型代表了ArbitraryType类型的指针类型。这里有4个与unsafe.Pointer类型相关的特殊转换操作。
1、一个指向其它类型值的指针值都可以被转换为一个unsafe.Pointer类型值。例如，如果有一个float32类型的变量f32，那么我们可以这样将与它对应的指针值
转换为一个unsafe.Pointer类型的值：pointer:=unsafe.Pointer(&f32)
2、一个unsafe.Pointer类型值可以被转换为一个与任何类型对应的指针类型的值。下面将pointer的值转换为与指向int类型值的指针值：
vptr:=(*int)(pointer)
3、一个unsafe.Pointer类型的值也可以被转换为一个uintptr类型的值。uptr:=uintptr(pointer)
4、一个uintptr类型的值也可以转换为一个unsafe.Pointer类型的值。pointer2：=unsafe.Pointer(uptr)
另外，我们还可以利用上述特殊转换操作以及unsafe包中声明的Offsetof函数进行有限的指针运算。
*/

type Person struct {
	Name    string `json:"name"`
	Age     uint8  `json:"age"`
	Address string `json:"addr"`
}

func main() {
	pp := &Person{"Robert", 32, "chengdu"}
	// 获取这个结构体值在内存中的存储地址
	var puptr = uintptr(unsafe.Pointer(pp))
	/**
	由于类型uintptr的值实际上是一个无符号整数，所以我们可以在该类型上进行任何算术运算，Offsetof函数会返回作为参数的某字段在其所属的结构体类型
	中的存储偏移量。换句话说，该函数的结果值就是在内存中从存储这个结构体值的起始位置到存储其中某字段的值的起始位置之间的距离。这个存储偏移量的单位是
	字节，它的类型是uintptr。实际上，同一个结构体类型的值在内存中的存储布局是固定的。也就是说，对于同一个结构体类型和它的同一个字段来说，这个存储
	偏移量总是相同的。我们知道了这个Person类型值的内存地址，也知道了它的存储起始位置到其中Name字段值的存储偏移量。依此，把它们相加就会得到存储这个
	结构体值中的Name字段值的内存地址。
	*/
	var npp uintptr = puptr + unsafe.Offsetof(pp.Name)
	// 获得了Name字段值的内存地址之后，可以利用上诉转换操作的第二条和第四条规则将它还原成指向这个Name字段值的指针类型值。
	var name *string = (*string)(unsafe.Pointer(npp))
	fmt.Println(*name) // Robert

	/**
	这里有个恒等式可以对上述示例中的一些操作进行很好的总结：
	uintptr(unsafe.Pointer(&s))+unsafe.Offsetof(s.f) == uintptr(unsafe.Pointer(&s.f))
	*/
}
