package main

import (
	"fmt"
	"sort"
)

/**
接口由方法集合代表，接口的类型声明由若干个方法的声明组成。
interface{}，空接口，它是不包含任何方法声明的接口。所有数据类型都是它的实现。
*/
type Interface interface {
	Len() int
	Less(i, j int) bool
	Swap(i, j int)
}

/**
可以将一个接口类型嵌入到另一个接口类型中，请记住：一个接口类型只接受其它接口类型的嵌入。不能嵌入自身，包括直接嵌入和简介嵌入。
*/
type Sortable interface {
	// 嵌入sort中的接口类型Interface
	sort.Interface
	Sort()
}

/**
实现接口
实现一个接口的类型可以是任何自定义的数据类型，只要这个数据类型附带的方法集合是该接口类型的方法集合的超级即可。
只要 funcName，args...，results...和interface一样，就认为这个数据类型实现了接口，不管方法的接受体是实例还是指针。
使用结构体类型来实现接口类型是更常用的一种做法。
*/
type SortableStrings [3]string

/**
把别名类型改为切片类型，这使得下面与索引表达式相关的方法不能通过编译。索引表达式不能被应用在指向切片值的指针类型值上。内建函数len的参数也不能是指向切片
值的指针类型值。解决办法很简单，即将方法Len、Less、Swap、和Sort的接收者类型由*SortableStrings改为SortableStrings。由于切片是引用类型，所以
值方法对接收者值的改变也会反应在其源值上。所以能正常排序。
*/
// type SortableStrings []string
func (self *SortableStrings) Len() int {
	return len(self)
}

func (self *SortableStrings) Less(i, j int) bool {
	return self[i] < self[j]
}

func (self *SortableStrings) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

/**
有了上面三个方法声明，SortableStrings类型就已经是一个sort.Interface接口类型的实现了。
*/

/**
一个接口类型可以被任意数量的数据类型实现。一个数据类型也可以实现多个接口了类型。
下面这个方法声明后表示SortableStrings实现了Sortable
*/
/*func (self SortableStrings) Sort() {
	sort.Sort(self)
}*/

/**
把方法接收者的类型改为SortableStrings类型对应的指针类型，这样SortableStrings类型就不再是接口类型Sortable的实现了。
因为此时的Sort方法不再是一个值方法了，而是一个指针方法了。所以下面的ok2是false。
这个时候只有与SortableStrings类型的值对应的指针值才能通过断言。
*/
func (self *SortableStrings) Sort() {
	sort.Sort(self)
}

type Fragment interface {
	Exec()
}

type FragmentImpl struct {
	i int
}

func (f *FragmentImpl) Exec() {
	fmt.Printf("%p\n", f)
	f.i = 10
}

type FragmentImp2 struct {
	i int
}

func (f FragmentImp2) Exec() {
	fmt.Printf("%p\n", f)
	f.i = 10
}

func main() {
	// _表示我们不关心类型转换后的结果，只关系转换成功与否。
	_, ok := interface{}(SortableStrings{}).(sort.Interface)
	_, ok1 := interface{}(SortableStrings{}).(Interface)
	_, ok2 := interface{}(SortableStrings{}).(Sortable)
	_, ok3 := interface{}(&SortableStrings{}).(Sortable)
	fmt.Println(ok)  // true，说明SortableStrings类型是sort.Interface接口类型的实现
	fmt.Println(ok1) // true，说明SortableStrings类型是Interface接口类型的实现
	fmt.Println(ok2) // false
	fmt.Println(ok3) // true

	ss := SortableStrings{"2", "3", "1"}
	ss.Sort()
	// 类型SortableStrings对应的指针类型的方法集合中包含了方法Sort，那么调用ss.Sort()即是调用(&ss).Sort()的速记法
	// (&ss).Sort()
	fmt.Printf("Sortable strings:%v\n", ss) // Sortable strings:[2 3 1]
	/**
	从上面的结果可以看出变量ss的值并没有被排序，SortableStrings类型的Sort方法实际上是通过sort.Sort来对接收者的值进行排序的。sort.Sort函数
	接受一个类型为sort.Interface的参数值，并利用这个值的方法Len、Less和Swap来修改其参数中的各个元素的位置以完成排序工作。SortableStrings声明
	的Len、Less和Swap方法都是值方法，在这些方法中对接收者值的改变并不会影响到它的源值。当把这些方法改为指针方法后就可以正常排序了。但是，此时的
	SortableStrings类型就已经不再是接口sort.Interface的实现了。所以上面的ok的值就是false。
	*/

	// 类型断言x.(T)，x表示的是接口类型，所以需要把123转换成接口类型
	value, ok4 := (interface{}(123)).(int)
	fmt.Println(value)
	fmt.Println(ok4) // true

	// 如果实现接口的是指针方法，那么不能把结构体实例值赋值给接口。
	//var f Fragment=FragmentImpl{1} // FragmentImpl does not implement Fragment (Exec method has pointer receiver)
	var f Fragment = &FragmentImpl{1}
	f.Exec()
	// 如果实现接口的是值方法，那么可以把结构体实例值或指针赋值给接口。
	var f2 Fragment = &FragmentImp2{1}
	f2.Exec()
	var f3 Fragment = FragmentImp2{1}
	f3.Exec()
}
