package order

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
)

type CompareFunction func(interface{}, interface{}) int8

/**
在Go语言中，具备有序性的预定义数据类型只有整数类型、浮点类型和字符串类型。
对字段keys的值有以下要求：
1、元素值应该都是有序的，我们应该可以方便地比较它们之间的大小。
2、元素类型不应该是一个具体的类型。我们应该可以在运行时再确定它的元素类型。
3、我们应该可以方便地对字段keys的值进行添加、删除以及获取等操作，就像对待一个普通的切片值那样。
4、元素值应该可以按照固定的顺序获取。
5、元素值应该能被自动排序。
6、由于字段keys的值总是已排序的，我们可以确定某一个元素值的具体位置。
7、既然我们可以在运行时决定字段keys的元素类型，那么也应该可以在运行时获知这个元素类型。
8、我们应该可以在运行时获取到被用于比较keys中不同元素值的大小的具体方法。
*/
type Keys interface {
	// 添加sort.Interface接口类型，就意味着Keys类型的值一定是可排序的
	sort.Interface
	Add(k interface{}) bool
	Remove(k interface{}) bool
	Clear()
	Get(index int) interface{}
	GetAll() []interface{}
	Search(k interface{}) (index int, contains bool)
	// reflect包中的程序实体为我们提供了Go语言运行时的反射机制。通过它们，我们可以编写一些代码来动态的操纵任意类型的对象
	ElementType() reflect.Type
	// 比较keys中不同元素值的大小的具体方法
	CompareFunc() CompareFunction
}

/**
为了能够动态地决定元素类型，我们不得不在Keys的实现类型中声明一个[]interface{}类型的字段，以作为存储被添加到Keys类型值中的元素的底层数据结构;
由于接口类型的值不具备有序性，不能比较它们的大小。不过，也许把这个问题抛出去并让使用Keys的实现类型的编程人员来解决它是一个可行的方案。因为他们应该
知道添加到Keys类型值中的元素值的实际类型并知道应该怎么比较它们；
*/
type myKeys struct {
	container []interface{}
	/**
	当第一个参数值小于第二个参数值时，结果值应该小于0；当第一个参数值大于第二个参数值时，结果值应该大于0；
	当第一个参数值等于第二个参数值时，结果值应该等于0
	*/
	compareFunc CompareFunction
	/**
	由于container字段是[]interface{}类型的，所以我们常常不能很方便地在运行时获取到它的实际元素类型。因此，我们需要一个明确container字段的实际
	元素类型的字段，这个字段的值所代表的类型也应该是当前Keys类型值的实际元素类型
	*/
	elementType reflect.Type
}

/**
实现sort.Interface接口中的几个方法
*/
func (keys *myKeys) Len() int {
	return len(keys.container)
}

func (keys *myKeys) Less(i, j int) bool {
	return keys.compareFunc(keys.container[i], keys.container[j]) == -1
}

func (keys *myKeys) Swap(i, j int) {
	keys.container[i], keys.container[j] = keys.container[j], keys.container[i]
}

/**
在真正向字段container添加元素值之前，应该先判断这个元素值的类型是否符合要求
*/
func (keys *myKeys) isAcceptableElement(k interface{}) bool {
	if k == nil {
		return false
	}
	// reflect.TypeOf函数确定参数k的实际类型
	if reflect.TypeOf(k) != keys.elementType {
		return false
	}
	return true
}

func (keys *myKeys) Add(k interface{}) bool {
	ok := keys.isAcceptableElement(k)
	if !ok {
		return false
	}
	keys.container = append(keys.container, k)
	sort.Sort(keys)
	return true
}

func (keys *myKeys) Search(k interface{}) (index int, contains bool) {
	ok := keys.isAcceptableElement(k)
	if !ok {
		return -1, false
	}
	/**
	sort.Search函数有两个参数。第一个参数接受的是要排序的切片值的长度，而第二个参数接受的是函数值。这个函数的意义是：对于一个给定的索引值，判定与之
	对应的元素值是否在要查找的元素值的右边。sort.Search函数使用二分查找算法，它要求被搜索的切片必须是有序的。
	*/
	index = sort.Search(keys.Len(), func(i int) bool {
		return keys.compareFunc(keys.container[i], k) >= 0
	})
	if index < keys.Len() && keys.container[index] == k {
		contains = true
	}
	return
}

func (keys *myKeys) Remove(k interface{}) bool {
	index, contains := keys.Search(k)
	if !contains {
		return false
	}
	keys.container = append(keys.container[:index], keys.container[index+1:]...)
	return true
}

func (keys myKeys) Clear() {
	keys.container = make([]interface{}, 0)
}

func (keys myKeys) Get(index int) interface{} {
	if index < 0 || index >= keys.Len() {
		return nil
	}
	return keys.container[index]
}

func (keys *myKeys) GetAll() []interface{} {
	initialLen := len(keys.container)
	snapshot := make([]interface{}, initialLen)
	actualLen := 0
	for _, key := range keys.container {
		if actualLen < initialLen {
			snapshot[actualLen] = key
		} else { // 如果在迭代完成之前，container的值中的元素有所增加，致使实际迭代次数大于initialLen，那么就使用append追加元素值
			snapshot = append(snapshot, key)
		}
		actualLen++
	}
	// 由于刚初始化snapshot时，它的元素值默认都是nil，如果在迭代完成之前，container的值中的元素有所减少，那么snapshot里面就有若干元素的值是nil，没哟意义
	// 所有要把这些nil去掉
	if actualLen < initialLen {
		snapshot = snapshot[:actualLen]
	}
	return snapshot
}

func (keys *myKeys) ElementType() reflect.Type {
	return keys.elementType
}

func (keys *myKeys) CompareFunc() CompareFunction {
	return keys.compareFunc
}

func (keys *myKeys) String() string {
	var buf bytes.Buffer
	buf.WriteString("keys<")
	buf.WriteString(keys.elementType.Kind().String())
	buf.WriteString(">{")
	first := true
	buf.WriteString("[")
	for _, key := range keys.container {
		if first {
			first = false
		} else {
			buf.WriteString(" ")
		}
		buf.WriteString(fmt.Sprintf("%v", key))
	}
	buf.WriteString("]")
	buf.WriteString("}")
	return buf.String()
}

func NewKeys(compareFunc CompareFunction, elementType reflect.Type) Keys {
	return &myKeys{
		container:   make([]interface{}, 0),
		compareFunc: compareFunc,
		elementType: elementType,
	}
}
