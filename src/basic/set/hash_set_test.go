package set

import (
	"runtime/debug"
	"testing"
)

/**
功能测试
测试源码文件总应该与被它测试的源码文件处在同一个代码包内，被用于测试程序实体功能的函数的名称和签名形如：
func TestXxx (t *testing.T)

测试覆盖率
作为被测试对象的代码包中的代码有多少在刚刚执行的测试中被使用到。换句话说，如果执行的该测试致使代码包中的80%的语句都被执行了，那么该测试的测试覆盖率就
是80%。go test命令可接受的与测试覆盖率有关的标记：
-cover 启用测试覆盖率分析
-covermode 自动添加-cover标记并设置不同的测试覆盖率统计模式。
set：只记录语句是否被执行过
count：记录语句被执行的次数
atomic：记录语句被执行的次数，并保证并发执行时也能正确计数，但性能会受到一定影响
这几个模式不可以被同时使用
-coverpkg 自动添加-cover标记并对该标记后罗列的代码包中的程序进行覆盖率统计。比如：-coverpkg bufio,net；在默认情况下，测试运行程序只会对被直接
测试的代码包中的程序进行统计。该标记意味着在测试中被间接使用到的其它代码包中的程序也可以被统计。
-coverprofile 自动添加-cover标记并把所有已通过的测试的覆盖率的概要写入指定的文件中

chenjunjiang@chenjunjiang-B85-HD3:~/go_workspace/goc2p/src$ go test -cover basic/set
ok  	basic/set	0.004s	coverage: 44.6% of statements
chenjunjiang@chenjunjiang-B85-HD3:~/go_workspace/goc2p/src$ go test -cover basic/set -coverprofile=cover.out
ok  	basic/set	0.003s	coverage: 44.6% of statements
chenjunjiang@chenjunjiang-B85-HD3:~/go_workspace/goc2p/src$ ls
basic  cover.out  hello  vendor
chenjunjiang@chenjunjiang-B85-HD3:~/go_workspace/goc2p/src$ go tool cover -func=cover.out
basic/set/hash_set.go:15:	NewHashSet	100.0%
basic/set/hash_set.go:25:	Add		75.0%
basic/set/hash_set.go:33:	Remove		100.0%
basic/set/hash_set.go:41:	Clear		100.0%
basic/set/hash_set.go:45:	Contains	100.0%
basic/set/hash_set.go:49:	Len		100.0%
basic/set/hash_set.go:53:	Same		0.0%
basic/set/hash_set.go:72:	Elements	81.8%
basic/set/hash_set.go:95:	String		100.0%
basic/set/hash_set.go:114:	IsSuperset	0.0%
basic/set/hash_set.go:137:	IsSuperset	0.0%
basic/set/set.go:14:		IsSet		66.7%
total:				(statements)	44.6%

chenjunjiang@chenjunjiang-B85-HD3:~/go_workspace/goc2p/src$ go tool cover -html=cover.out
执行之后立马会在浏览器打开生成的html文件，在这个页面，被测试到的语句以绿色显示，未被测试到的语句以红色显示，而未参加测试覆盖率计算的语句则用灰色表示。
这个HTML文件展示的是在set统计模式下生成的概要文件。
*/
func TestNewHashSet(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			// 当被测试的程序实体的状态不正确的时候，调用t.Errorf方法及时对当前的错误状态进行记录，它相当于先后对t.Logf和t.Fail方法进行调用
			t.Errorf("Fatal Error:%s\n", err)
		}
	}()
	t.Log("Starting TestNewHashSet......")
	hs := NewHashSet()
	t.Logf("Create a HashSet value: %v\n", hs)
	if hs == nil {
		t.Errorf("The result of func NewHashSet is nil!\n")
	}
	isSet := IsSet(hs)
	if !isSet {
		t.Errorf("The value of HashSet is not Set!\n")
	} else {
		t.Logf("The HashSet value is a Set.\n")
	}
}

func TestHashSet_Len(t *testing.T) {
	testSetLenAndContains(t, func() Set {
		return NewHashSet()
	}, "HashSet")
}

func TestHashSet_Contains(t *testing.T) {
	testSetLenAndContains(t, func() Set {
		return NewHashSet()
	}, "HashSet")
}

func TestHashSet_Add(t *testing.T) {
	testSetAdd(t, func() Set {
		return NewHashSet()
	}, "HashSet")
}

func TestHashSet_Remove(t *testing.T) {
	testSetRemove(t, func() Set {
		return NewHashSet()
	}, "HashSet")
}

func TestHashSet_Clear(t *testing.T) {
	testSetClear(t, func() Set {
		return NewHashSet()
	}, "HashSet")
}

func TestHashSet_Elements(t *testing.T) {
	testSetElements(t, func() Set {
		return NewHashSet()
	}, "HashSet")
}

func TestHashSet_String(t *testing.T) {
	testSetString(t, func() Set {
		return NewHashSet()
	}, "HashSet")
}
