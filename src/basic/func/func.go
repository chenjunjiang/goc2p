package main

import "fmt"

func Module(x, y int) (result int) {
	return x % y
}

type MyInt int

func (i MyInt) Test() {

}

func main() {
	fmt.Println(Module(7, 3))
}
