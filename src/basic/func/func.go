package main

import "fmt"

func Module(x, y int) (result int) {
	return x % y
}

func Module1(x, y int) int {
	return x % y
}

/**
函数在被调用的时候就会初始化对应结果类型变量的零值（这里result就初始化为0）,如果这样的函数的函数体内有一条不带任何参数的return语句，那么在执行
return语句的时候，作为结果的变量的当前值就会返回给函数调用方。
*/
func Module2(x, y int) (result int) {
	return
}

/**
result的值作为函数调用方的返回值
*/
func Module3(x, y int) (result int) {
	result = x + y
	return
}

/**
定义一个匿名函数
*/
var Module4 = func(x, y int) (result int) {
	result = x + y
	return result
}

/**
函数既可以作为变量的值，也可以做为其它函数的参数或返回结果
*/
// 声明一个函数类型,type关键字专门用于声明自定义数据类型
type Encipher func(plaintext string) []byte

func GenEncryptionFunc(encipher Encipher) func(string) (ciphertext string) {
	return func(plaintext string) string {
		return fmt.Sprintf("%x", encipher(plaintext))
	}
}

/**
方法就是附属于某个自定义的数据类型的函数。具体来说，一个方法就是一个与某个接收者关联的函数。方法的声明包含了关键字func、接收者声明、方法名称、参数声明
列表、结果声明列表和方法体。
方法格式:
func (o Object) FuncName(args...) (results...){}
其中 o 代表的是函数的接受体，意思是这个函数属于对象 o 的方法,args 表示形参列表,results 表示函数返回值列表，对于无返回值的方法可以为空
接收者声明有关的几条编写规则：
1、接收者声明中的类型必须是某个自定义的数据类型，或者是一个与某自定义数据类型对于的指针类型。但不论接收者的类型是哪一种，接收者的基本类型都会是
那个自定义的数据类型。接收者的基本类型基本既不能是一个指针类型，也不能是一个接口类型。
2、接收者声明中的类型必须由非限定标识符代表。也就是说，方法所属的数据类型的声明必须与该方法声明处在同一个代码包内。
3、接收者标识符不能是空标识符"_"，并且必须在其所在的方法签名中是唯一的。
4、如果接收者的值未在当前方法的方法体内被引用，那我们就可以将这个接收者标识符从当前方法的接收者声明中删除掉。但并不推荐这么做。
我们通常把接收者类型是某个自定义数据类型的方法叫做该数据类型的值方法，而把接收者类型是与某个自定义数据类型对应的指针类型的方法叫做指针方法。
一个方法的类型与从其声明中去掉接收者声明之后的函数的类型相似。
func (self *MyIntSlice) Min() (result int)的类型是func Min() (self *MyIntSlice, result int)
*/
type MyIntSlice []int

func (i MyIntSlice) Max() (result int) {
	return
}

type FragmentImpl struct {
	i int
}

func (f FragmentImpl) Exec() {
	fmt.Printf("%p\n", &f)
	f.i = 10
}
func (f *FragmentImpl) Exec2() {
	fmt.Printf("%p\n", f)
	f.i = 10
}

func main() {
	fmt.Println(Module(7, 3))
	fmt.Println(Module2(7, 3)) // 0
	fmt.Println(Module3(7, 3)) // 10
	fmt.Println(Module4(7, 3)) // 10
	myIntSlice := MyIntSlice{1, 2}
	fmt.Println(myIntSlice.Max()) // 0

	/*encipher := func(plaintext string) []byte {
		return []byte{1, 2}
	}
	GenEncryptionFunc(encipher)*/

	fragment := FragmentImpl{1}
	fmt.Printf("%p --  %v \n", &fragment, fragment)
	fragment.Exec()
	fmt.Printf("%p --  %v \n", &fragment, fragment)
	(&fragment).Exec()
	fmt.Printf("%p --  %v \n", &fragment, fragment)

	fmt.Println("----------------------------")

	fragment2 := &FragmentImpl{1}
	fmt.Printf("%p --  %v \n", fragment2, fragment2)
	fragment2.Exec2()
	fmt.Printf("%p --  %v \n", fragment2, fragment2)
	(*fragment2).Exec2()
	fmt.Printf("%p --  %v \n", fragment2, fragment2)

	/**
	可以看出，对于Exec()，虽然方法接受体是实例，但是不管FragmentImpl的实例还是实例的地址指针都能够调用成功，但是不管以何种形式调用，
	都会产生值复制（地址不同），不会对原始对象的属性产生更改。

	对于Exec2()，方法的接受体是结构体指针，同样的也可以用实例或者指针去调用，而且能够对原始对象产生修改行为。

	所以，在golong中，会不会产生值复制，关键在于方法的接受体是对象还是指针，而不是调用方法的类型，因为golang会自动解引用。
	不过我觉得，即使以*T指针的形式调用方法，实际上也是会产生复制的，不过此时复制的是指针这个“值”，所以在方法中可以通过访问指针修改原生的值。
	*/
}
