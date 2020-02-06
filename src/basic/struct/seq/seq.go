package seq

import (
	"fmt"
	"reflect"
	"sort"
)

type Sortable interface {
	sort.Interface
	Sort()
}

/**
为了使Sequence类型能够部分模拟泛型类型的行为特征，只向它嵌入Sortable接口类型是不够的。我们需要对Sortable接口类型进行扩展（为了符合CP原则），而不是
修改接口。我们需要创建一个新接口类型，并将Sortable接口类型嵌入其中。
一个类型想要实现GenericSeq，也必须要实现Sortable
*/
type GenericSeq interface {
	Sortable
	Append(e interface{}) bool
	Set(index int, e interface{}) bool
	Delete(index int) (interface{}, bool)
	ElemValue(index int) interface{}
	ElemType() reflect.Type
	Value() interface{}
}

/**
elemType用来缓存GenericSeq字段中存储的值的元素类型。
*/
type Sequence struct {
	GenericSeq
	Sorted1   bool
	ElemType1 reflect.Type
}

/**
为了能够在改变GenericSeq字段存储的值的过程中及时对字段sorted和elemType的值就行修改，我们还创建了几个与Sequence类型相关联的方法。
这些方法分别与接口类型GenericSeq或Sortable中声明的某个方法有着相同的方法名称和方法签名。也就是说，我们通过这种方式隐藏了GenericSeq字段中存储的
值的这些同名方法。
*/
func (self *Sequence) Sort() {
	self.GenericSeq.Sort()
	self.Sorted1 = true
}

func (self *Sequence) Append(e interface{}) bool {
	result := self.GenericSeq.Append(e)
	if result && self.Sorted1 {
		self.Sorted1 = false
	}
	return result
}

func (self *Sequence) Set(index int, e interface{}) bool {
	result := self.GenericSeq.Set(index, e)
	if result && self.Sorted1 {
		self.Sorted1 = false
	}
	return result
}

func (self *Sequence) ElemType() reflect.Type {
	if self.ElemType1 == nil && self.GenericSeq != nil {
		self.ElemType1 = self.GenericSeq.ElemType()
	}
	return self.ElemType1
}

func (self *Sequence) Init() (ok bool) {
	if self.GenericSeq != nil {
		self.ElemType1 = self.GenericSeq.ElemType()
		ok = true
	}
	return ok
}

func (self *Sequence) Sorted() bool {
	return self.Sorted1
}

/**
在初始化Sequence类型值的时候，我们还需要用到GenericSeq接口类型的实现类型。下面是这个实现
*/
type StringSeq struct {
	Str []string // 这里字段名称首字母大写的目的是为了让其它包能使用
}

func (self *StringSeq) Len() int {
	return len(self.Str)
}

func (self *StringSeq) Less(i, j int) bool {
	return self.Str[i] < self.Str[j]
}

func (self *StringSeq) Swap(i, j int) {
	self.Str[i], self.Str[j] = self.Str[j], self.Str[i]
}

func (self *StringSeq) Sort() {
	sort.Sort(self)
}

func (self *StringSeq) Append(e interface{}) bool {
	s, ok := e.(string)
	if !ok {
		return false
	}
	self.Str = append(self.Str, s)
	return true
}

func (self *StringSeq) Set(index int, e interface{}) bool {
	if index >= self.Len() {
		return false
	}
	s, ok := e.(string)
	if !ok {
		return false
	}
	self.Str[index] = s
	return true
}

func (self *StringSeq) Delete(index int) (interface{}, bool) {
	length := self.Len()
	if index >= length {
		return nil, false
	}
	s := self.Str[index]
	if index < (length - 1) {
		copy(self.Str[index:], self.Str[index+1:])
	}
	invalidIndex := length - 1
	self.Str[invalidIndex] = ""
	self.Str = self.Str[:invalidIndex]
	return s, true
}

func (self StringSeq) ElemValue(index int) interface{} {
	if index >= self.Len() {
		return nil
	}
	return self.Str[index]
}

func (self *StringSeq) ElemType() reflect.Type {
	return reflect.TypeOf(self.Str).Elem()
}

func (self StringSeq) Value() interface{} {
	return self.Str
}

func MyPack() {
	fmt.Println("这是我的第一自定义包")
}
