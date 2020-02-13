package set

import (
	"bytes"
	"fmt"
)

type HashSet struct {
	m map[interface{}]bool
}

/**
结果声明的类型是*HashSet而不是HashSet，是因为希望在这个结果值的方法集合中包含调用接收者类型为HashSet或*HashSet的所有方法
*/
func NewHashSet() *HashSet {
	return &HashSet{m: make(map[interface{}]bool)}
}

/**
这里接收者的类型是*HashSet而不是HashSet，主要原因是减少复制接收者值时对系统资源的耗费。方法的接收者值只是当前值的一个复制品。所以，当Add方法的
接收者的类型为HashSet的时候，对它的每一次调用都需要对当前值(HashSet类型值)进行一次复制。当Add方法的接收者类型为*HashSet时，对它进行调用时复制
的当前值(*HashSet类型值)只是一个指针值。在大多数情况下，一个指针值占用的内存空间总会比它指向的那个其它类型的值所占的内存空间小。因此，从节约内存
空间的角度出发，建议尽量将方法的接收者类型设置为相应的指针类型。
*/
func (set *HashSet) Add(e interface{}) bool {
	if !set.m[e] {
		set.m[e] = true
		return true
	}
	return false
}

func (set *HashSet) Remove(e interface{}) {
	delete(set.m, e)
}

/**
这里不能通过遍历的方式一个个删除，在并发情况下会出问题。注意:这里必须是*HashSet，如果用HashSet，那么只是为当前值的某个复制品的字段m赋值而已。
已经与字段m解除绑定的那个旧字典值由于不再与任何程序实体存在绑定关系而成为无用的数据，它会在之后的某一时刻被Go语言的垃圾回收器发现并回收。
*/
func (set *HashSet) Clear() {
	set.m = make(map[interface{}]bool)
}

func (set *HashSet) Contains(e interface{}) bool {
	return set.m[e]
}

func (set *HashSet) Len() int {
	return len(set.m)
}

func (set *HashSet) Same(other Set) bool {
	if other == nil {
		return false
	}
	if set.Len() != other.Len() {
		return false
	}
	for key := range set.m {
		if !other.Contains(key) {
			return false
		}
	}
	return true
}

/**
迭代Set，虽然里面通过两个判断来处理在迭代过程中m的值会有所增加或减少的情况， 但是在并发情况下还是不能保证Elements总是执行正确，要做到真正的并发安全
，需要用到读写互斥量。
*/
func (set *HashSet) Elements() []interface{} {
	initialLen := len(set.m)
	snapshot := make([]interface{}, initialLen)
	actualLen := 0
	for key := range set.m {
		if actualLen < initialLen {
			snapshot[actualLen] = key
		} else { // 如果在迭代完成之前，m的值中的元素有所增加，致使实际迭代次数大于initialLen，那么就使用append追加元素值
			snapshot = append(snapshot, key)
		}
		actualLen++
	}
	// 由于刚初始化snapshot时，它的元素值默认都是nil，如果在迭代完成之前，m的值中的元素有所减少，那么snapshot里面就有若干元素的值是nil，没哟意义
	// 所有要把这些nil去掉
	if actualLen < initialLen {
		snapshot = snapshot[:actualLen]
	}
	return snapshot
}

/**
获取Set的字符串表现形式，使用bytes.Buffer类型值作为结果值的缓冲区，这样可以避免因string类型值的拼接造成内存空间上的浪费
*/
func (set *HashSet) String() string {
	var buf bytes.Buffer
	buf.WriteString("Set{")
	first := true
	for key := range set.m {
		if first {
			first = false
		} else {
			buf.WriteString(" ")
		}
		buf.WriteString(fmt.Sprintf("%v", key))
	}
	buf.WriteString("}")
	return buf.String()
}

/**
A集合是否真包含B，也就是说集合A是不是B的一个超集
*/
func (one *HashSet) IsSuperset(other *HashSet) bool {
	if other == nil {
		return false
	}
	oneLen := one.Len()
	otherLen := other.Len()
	if oneLen == 0 || oneLen == otherLen {
		return false
	}
	if oneLen > 0 && otherLen == 0 {
		return true
	}
	for _, v := range other.Elements() {
		if !one.Contains(v) {
			return false
		}
	}
	return true
}

/**
把IsSuperset方法抽离出来，并使之成为独立的函数，以面向所有的实现类型
*/
func IsSuperset(one Set, other Set) bool {
	if other == nil {
		return false
	}
	oneLen := one.Len()
	otherLen := other.Len()
	if oneLen == 0 || oneLen == otherLen {
		return false
	}
	if oneLen > 0 && otherLen == 0 {
		return true
	}
	for _, v := range other.Elements() {
		if !one.Contains(v) {
			return false
		}
	}
	return true
}
