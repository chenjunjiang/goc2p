package main

import (
	. "basic/map/order"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

/**
map[K]T，字典声明中的元素类型可以是任意一个有效的Go语言数据类型，但是，它的键类型不能是函数类型、字典类型或切片类型。如果键类型是接口类型，那么就
要求在程序运行期间，该类型的字典值中的每一个键值的动态类型都必须是可比较的(用操作符==)。
*/
func main() {
	editorSign := map[string]bool{"Vim": true, "LiteIDE": true, "Notepad": false}
	fmt.Println(len(editorSign)) // 3
	sign, ok := editorSign["Vim"]
	fmt.Println(sign)
	fmt.Println(ok)
	// strconv.FormatBool把bool类型转换为string类型
	// string类型和其他类型的值的互转 (https://blog.csdn.net/bobodem/article/details/80182096)
	fmt.Println(strings.Join([]string{strconv.FormatBool(sign), strconv.FormatBool(ok)}, ","))
	editorSign["Vim"] = false
	sign, ok = editorSign["Vim"]
	fmt.Println(sign) // false
	fmt.Println(ok)
	delete(editorSign, "Vim")
	fmt.Println(len(editorSign)) // 2

	fmt.Println("测试有序Map......")
	intKeys := NewKeys(func(i1 interface{}, i2 interface{}) int8 {
		// 把i1转换成int
		k1 := i1.(int)
		k2 := i2.(int)
		if k1 < k2 {
			return -1
		} else if k1 > k2 {
			return 1
		} else {
			return 0
		}
	}, reflect.TypeOf(int(1)))
	orderMap := NewOrderedMap(intKeys, reflect.TypeOf(""))
	orderMap.Put(3, "test3")
	orderMap.Put(1, "test1")
	orderMap.Put(5, "test5")
	orderMap.Put(2, "test2")
	orderMap.Put(9, "test9")
	orderMap.Put(8, "test8")
	orderMap.Put(6, "test6")
	fmt.Println(orderMap.String())
	fmt.Println(intKeys.Search(2))
	fmt.Println(orderMap.HeadMap(5).String())
}
