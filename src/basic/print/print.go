package main

import "fmt"

func main() {
	var strs []interface{} = []interface{}{"a", "b", "c"}
	// r:=fmt.Sprint(strs) //[a b c]
	r := fmt.Sprint(strs...) // abc, ...的意思是把strs打散
	r1 := fmt.Sprint("a", "b", "c")
	fmt.Println(r)
	fmt.Println(r1) // abc
	output("a", "b", "c")
}

func output(v ...interface{}) {
	r := fmt.Sprint(v...)
	fmt.Println(r) // abc
	r = fmt.Sprint(v)
	fmt.Println(r) // [a b c]
}
