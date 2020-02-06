package main

import "fmt"

func main() {
	var ipv4 [4]uint8 = [4]uint8{192, 168, 0, 1}
	// var ipv4 = [4]uint8{192, 168, 0, 1} // 类型可以省略
	// ipv4 := [4]uint8{192, 168, 0, 1}
	// ...表示需由Go编译器计算该值的元素数量并以此获得其长度
	// ipv4 := [...]uint8{192, 168, 0, 1}
	fmt.Println(len(ipv4)) // 长度 4
	fmt.Println(cap(ipv4)) // 容量 4
	array := [2]uint8{}
	array[0] = 1
	array[1] = 2
	fmt.Println(array)
}
