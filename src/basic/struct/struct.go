package main

import (
	. "basic/set"
	"basic/struct/seq"
	"fmt"
	"sort"
)

/**
匿名字段的类型必须由一个数据类型的名称或者一个与非接口类型对应的指针类型的名称代表。代表匿名字段类型的非限定名称将被隐含地作为该字段的名称。如果匿名
字段是一个指针类型的话，那么这个指针类型所指的数据类型的非限定名称就会作为该字段的名称，所谓非限定名称就是由非限定标识符代表的名称，非限定标识符指的是
不包含代码包名称和点"."的标识符。
type A struct{
    T1
    *T2
    P.T3
    *P.T4
}
T1和P.T3隐含的名称是T1和T3，*T2和*P.T4隐含的名称是T2和T4
*/
/*type Sequence struct {
	len           int
	cap           int
	Sortable      // 匿名字段
	sortableArray sort.Interface
}*/
type Sequence struct {
	Sortable // 匿名字段
	sorted   bool
}

type Sortable interface {
	Sort()
	Sort1()
}

type SortableStrings []string

func (self SortableStrings) Len() int {
	return len(self)
}

func (self SortableStrings) Less(i, j int) bool {
	return self[i] < self[j]
}

func (self SortableStrings) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

func (self SortableStrings) Sort() {
	fmt.Println("Sort被调用了......")
	sort.Sort(self)
}

func (self *SortableStrings) Sort1() {
	fmt.Println("Sort1被调用了......")
	sort.Sort(self)
}

func (self Sequence) test() {
	fmt.Println("Sequence test")
}

func (self *Sequence) test1() {
	fmt.Println("Sequence test1")
}

/**
假如Sequence类型中包含了一个与Sortable接口类型的Sort方法的名称和签名都相同的方法，那么seq.Sort()就一定是对Sequence类型值自身附带的Sort方法
的调用。也就是说，在这种情况下，嵌入类型Sortable的Sort方法被隐藏了。
*/
func (self *Sequence) Sort() {
	fmt.Println("Sequence自己的Sort方法被调用了")
	self.Sortable.Sort()
	self.sorted = true
}

/**
在字段声明的后面添加一个字符串字面量标签，以作为对应字段的附加属性。这种标签对于使用使用该结构体类型及其值的代码是不可见的，但是，我们可以用标准代码库
代码包reflect中提供的函数查看到结构体类型中字段的标签。因此这种标签常常会在一些特殊应用场景下使用，比如，标准库代码包encoding/json中的函数会根据
这种标签的内容确定与该结构体类型中的字段对应的JSON节点的名称。
*/
type people struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"addr"`
}

func (p *people) SetName(name string) {
	p.Name = name
}

func (p people) SetAge(age int) {
	p.Age = age
}

/**
在一个结构体类型的别名类型值上，既不能调用那个结构体类型的方法，也不调用与那个结构体类型对应的指针方法。如下面main方法里面的list.test()和list.test1()。
别名类型不是源类型的子类型。但是，别名类型内部的结构会与它的源类型一致。对于一个结构体类型的别名类型来说，它拥有源类型的全部字段。
正因为别名类型存在这样的局限性，我们才会在Sequence类型中嵌入了接口类型Sortable，而不是直接将Sequence类型声明为一个接口类型Sortable的某个实现类型。
只不过这样做的原因不只上面这一个。比如，嵌入字段Sortable能够用于存储所有实现了该接口类型的数据类型的值。这样的类型结构设计使得Sequence类型可以在一定
程度上模拟出泛型类型的一些特点，这个特点实现请参考seq.go。
*/
type List Sequence

func main() {
	ss := SortableStrings{"2", "3", "1"}
	// 如果实现接口的方法中有指针方法，那么不能把接收者(SortableStrings)的实例值赋值给接口(Sortable)。
	// seq := Sequence{ss, false}
	seq1 := Sequence{&ss, false}
	seq1.Sort()
	seq1.Sort1()
	fmt.Printf("Sortable strings:%v\n", ss) //Sortable strings:[1 2 3]
	fmt.Println(seq1.sorted)                // true

	/**
	不管方法的 receiver(接收者) 是对象的值还是指针，对象的值和指针均可以调用该方法。即对象的值既可以调用 receiver 是值的方法，也可以调用 receiver 是指针的方法。
	对象的指针也是如此。
	*/
	p := people{"zhangsan", 22, "beijing"}
	// p:=&people{"zhangsan",22}
	p.SetAge(11)
	p.SetName("lisi")

	/**
	匿名结构体与命名结构体相比，更像是"一次性"的类型。它不具有通用性，因此常常用在临时数据存储和传递的场景中。
	*/
	anonym := struct {
		a int
		b int
	}{1, 2}
	fmt.Println(anonym.a)
	fmt.Println(anonym.b)

	list := List{&SortableStrings{"2", "3", "1"}, false}
	fmt.Println(list.sorted)
	seq1.test()
	//list.test()
	//list.test1()

	// 调用seq包中的方法，在使用被导入代码包中的程序实体时，需要使用包路径的最后一个元素加"."的方式
	seq.MyPack()

	str := []string{"4", "7", "5"}
	// stringSeq:=seq.StringSeq{str} // Unnamed field initialization
	stringSeq := seq.StringSeq{Str: str} // 建议是加上结构体字段的名称
	// 如果实现接口的是指针方法，那么不能把结构体实例值赋值给接口。
	// seq2:=seq.Sequence{stringSeq,false,nil}
	seq2 := seq.Sequence{&stringSeq, false, nil}
	result := seq2.Append("8")
	fmt.Println(result)
	result = seq2.Set(2, "9")
	fmt.Println(result)
	seq2.Sort()
	fmt.Println(stringSeq.Str)

	fmt.Println("测试Set......")
	// 在使用被导入代码包中的程序实体时，如果不想加前缀，可以在导入代码包的时候使用"."来代替别名
	var set1 Set = NewHashSet()
	var set2 Set = NewHashSet()
	set1.Add(1)
	set1.Add(2)
	set1.Add(3)
	set2.Add("2")
	set2.Add("3")
	set2.Add("4")
	fmt.Println(set1.Same(set2))
	fmt.Println(set1.String())
	fmt.Println(IsSuperset(set1, set2))
}
