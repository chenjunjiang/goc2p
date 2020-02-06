package main

func main() {
	// \n会被转义
	var str string = "xxxx\nyyy"
	println(str)
	//\n不会被转义，原样输出
	str1 := `xxx\n`
	println(str1)
}
