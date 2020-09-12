package main

import "fmt"

/**
切片的长度是切片中元素的数量，能被访问的元素数量不能超过长度
切片的容量是从创建切片的索引开始到底层数组末尾的元素数量
切片是可索引的，并且可以由len()方法获取长度，切片提供了计算容量的方法cap()，可以测量切片最长可以达到多少。
切片实际的是获取数组的某一部分，len(切片)<=cap(切片)<=len(数组)
*/
func main() {
	// 切片可以看做一种对数组的包装形式，它包装的数组称为该切片的底层数组。反过来讲，切片是针对其底层数组中某个连续片段的描述。
	var ips = []string{"192.168.0.1", "192.168.0.2", "192.168.0.3"}
	fmt.Println(ips[0])
	/**
	与数组不同，切片的类型字面量(如[]string)并不携带长度信息。切片的长度是可以变的，且并不是类型的一部分；只要元素类型相同，两个切片的类型就是
	相同的。此外，一个切片类型的零值总是nil，此零值长度长度和容量都为0。
	切片值相当于堆某个底层数组的引用。其内部结构包含了3个元素：指向底层数组中某个元素的指针、切片的长度以及切片的容量。这里所说的容量是指，从指针指向
	的那个元素到底层数组的最后一个元素的元素个数。
	*/
	// 切片表达式
	//Go 语言对 slice 有两种表示方式：简略表达式与完整表达式。
	//Slice 的简略表达式是：Input[low:high],其中，low 和 high 是 slice 的索引(index)，其数值必须是整数，它们指定了输入操作数(Input)的哪些元素可以被放置在结果的 slice 中。
	//输入操作数可以是 string，array，或者是指向 array 或 slice 的指针。结果 slice 的长度就是 high-low。
	numbers1 := [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	s1 := numbers1[1:3]
	fmt.Println(s1)      // [1 2]
	fmt.Println(len(s1)) // 2
	fmt.Println(cap(s1)) // 9
	// Slice 的索引 low 和 high 可以省略，low 的默认值是0，high 的默认值为 slice 的长度：
	fmt.Println("foo"[:2])        // "fo"
	fmt.Println("foo"[1:])        // "oo"
	fmt.Println("foo"[:])         // "foo"
	fmt.Println((&numbers1)[1:3]) // [1, 2]
	fmt.Println("----------------")

	// 完整的 slice 表达式具有以下的形式：input[low:high:max]，max 将slice 的容量设置为 max-low。
	// 这种方法可以控制结果 slice 的容量，但是只能用于 array 和指向 array 或 slice 的指针（ string 不支持）
	numbers2 := [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	s2 := numbers2[2:4:6]
	fmt.Println(s2)      //[2 3]
	fmt.Println(len(s2)) // 2
	fmt.Println(cap(s2)) // 4
	fmt.Println(s2[1])
	// fmt.Println(s2[2]) // panic: runtime error: index out of range [2] with length 2
	// 扩展s2，这样就能看到最多的底层数组的元素值了。注意：一个切片的容量是固定的。也就是说，我能能够看到的底层数组元素的最大数量是固定的。
	s2 = s2[:cap(s2)]
	fmt.Println(s2) // [2 3 4 5]
	s2[2] = 22
	fmt.Println(s2[2]) // 22

	// 当 slice 的输入操作数是一个 slice 时，结果 slice 的容量取决于输入操作数，而不是它的指向的底层 array
	numbers3 := [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Println(cap(numbers3)) // 10
	s3 := numbers3[1:4]
	fmt.Println(s3)      // [1, 2, 3]
	fmt.Println(cap(s3)) // 9
	s4 := numbers3[1:4:5]
	fmt.Println(s4)      // [1, 2, 3]
	fmt.Println(cap(s4)) // 4
	s5 := s4[:]
	fmt.Println(s5)      // [1, 2, 3]
	fmt.Println(cap(s5)) // 4

	// append函数会将指定的若干元素追加到原切片值的末端。
	ips = append(ips, "192.168.0.4")
	println(ips[3])
	// 赋值操作会把这个新的切片值再赋给ips。注意，新、旧切片值可能指向不同的底层数组。若新切片值的底层数组的长度不足以完成元素的追加操作，它将会被
	// 更长的底层数组替换，以容纳更多的元素。

	// 内建函数make用于初始化切片、字典或通道类型的值。对于切片类型来说，用make函数的好处是可以用很短的代码初始化一个长度很大的值
	ips = make([]string, 100)
	fmt.Println(len(ips)) // 100
	fmt.Println(cap(ips)) // 100
	// 用make初始化的切片值的每一个元素值都会是其元素类型的零值，这里ips中的那100个元素的都会是空字符串""。
	fmt.Println(ips[0])  // ""
	fmt.Println(ips[99]) // ""

	// 创建长度为0，容量为10的切片
	s := make([]int, 0, 10)
	// 长度为0说明还切片中还不存在元素，访问会报错
	// fmt.Println(s[0])
	s = append(s, 1)
	fmt.Println(s[0])

	x := make([]int, 2, 10)
	fmt.Println(x[0]) // 0
	fmt.Println(x[1]) // 0
	// fmt.Println(x[2]) // index out of range [2] with length 2

	v := []int{2, 3, 4, 5, 6, 7}
	v = v[:0]
	fmt.Println(len(v)) // 0
	fmt.Println(cap(v)) // 6
	// fmt.Println(v[0]) // index out of range [0] with length 0
}
