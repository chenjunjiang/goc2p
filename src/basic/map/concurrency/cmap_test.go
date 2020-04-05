package concurrency

import (
	"bytes"
	"fmt"
	"math/rand"
	"reflect"
	"runtime/debug"
	"testing"
	"time"
)

func TestInt64Cmap(t *testing.T) {
	newCmap := func() ConcurrentMap {
		keyType := reflect.TypeOf(int64(2))
		elemType := keyType
		return NewConcurrentMap(keyType, elemType)
	}
	testConcurrentMap(
		t,
		newCmap,
		func() interface{} { return rand.Int63n(1000) },
		func() interface{} { return rand.Int63n(1000) },
		reflect.Int64,
		reflect.Int64)
}

func TestFloat64Cmap(t *testing.T) {
	newCmap := func() ConcurrentMap {
		keyType := reflect.TypeOf(float64(2))
		elemType := keyType
		return NewConcurrentMap(keyType, elemType)
	}
	testConcurrentMap(
		t,
		newCmap,
		func() interface{} { return rand.Float64() },
		func() interface{} { return rand.Float64() },
		reflect.Float64,
		reflect.Float64)
}

func TestStringCmap(t *testing.T) {
	newCmap := func() ConcurrentMap {
		keyType := reflect.TypeOf(string(2))
		elemType := keyType
		return NewConcurrentMap(keyType, elemType)
	}
	testConcurrentMap(
		t,
		newCmap,
		func() interface{} { return genRandString() },
		func() interface{} { return genRandString() },
		reflect.String,
		reflect.String)
}

/**
我们使用ResetTimer、StartTimer和StopTimer方法忽略掉一些无关紧要的语句的执行时间。我们只对
imap[key] = elem、_ = imap[key]和ml := len(imap)语句的执行时间计时
*/
func BenchmarkMap(b *testing.B) {
	keyType := reflect.TypeOf(int32(2))
	elemType := keyType
	imap := make(map[interface{}]interface{})
	var key, elem int32
	fmt.Printf("N=%d.\n", b.N)
	// ResetTimer是重置计时器，这样可以避免for循环之前的初始化代码的干扰
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 停止计时器
		b.StopTimer()
		seed := int32(i)
		key = seed
		elem = seed << 10
		// 开始计时器
		b.StartTimer()
		imap[key] = elem
		_ = imap[key]
		// 停止计时器
		b.StopTimer()
		// 记录单次操作中被处理的字节的数量，在这里，我们每次记录一个键值对所用的字节数。由于键和值类型都是int32，所以它们共会用掉8个字节
		b.SetBytes(8)
		b.StartTimer()
	}
	ml := len(imap)
	b.StopTimer()
	mapType := fmt.Sprintf("Map<%s, %s>",
		keyType.Kind().String(), elemType.Kind().String())
	b.Logf("The length of %s value is %d.\n", mapType, ml)
	b.StartTimer()
}

/**
执行命令后要等待一会儿才会出结果
chenjunjiang@chenjunjiang-B85-HD3:~/go_workspace/goc2p/src/basic/map/concurrency$ go test -bench="." -run="^$" -benchtime=1s -v
N=1.
goos: linux
goarch: amd64
pkg: basic/map/concurrency
BenchmarkMap-4                  N=100.
N=10000.
N=1000000.
N=1486742.
 1486742               718 ns/op          11.14 MB/s
--- BENCH: BenchmarkMap-4
    cmap_test.go:90: The length of Map<int32, int32> value is 1.
    cmap_test.go:90: The length of Map<int32, int32> value is 100.
    cmap_test.go:90: The length of Map<int32, int32> value is 10000.
    cmap_test.go:90: The length of Map<int32, int32> value is 1000000.
    cmap_test.go:90: The length of Map<int32, int32> value is 1486742.
N=1.
BenchmarkConcurrentMap-4        N=100.
N=10000.
N=1000000.
 1000000              1011 ns/op           7.91 MB/s
--- BENCH: BenchmarkConcurrentMap-4
    cmap_test.go:117: The length of ConcurrentMap<int32, int32> value is 1.
    cmap_test.go:117: The length of ConcurrentMap<int32, int32> value is 100.
    cmap_test.go:117: The length of ConcurrentMap<int32, int32> value is 10000.
    cmap_test.go:117: The length of ConcurrentMap<int32, int32> value is 1000000.
PASS
ok      basic/map/concurrency   105.152s

BenchmarkConcurrentMap在1秒内最多执行相关操作的次数约为1000000次，平均每次执行的耗时是1011纳秒，每秒处理7.91 MB的数据
BenchmarkMap在1秒内最多执行相关操作的次数约为1486742次，平均每次执行的耗时是718纳秒，每秒处理11.14 MB的数据
这种差距主要是ConcurrentMap中使用读写锁造成的。
同步工具在为程序并发安全提供支持的同时也会对其性能造成不可忽略的损耗。还有Go并未对自定义泛型提供支持，以至于我们编写此类扩展的时候并不是那么方便，
有的时候，我们不得不使用反射API，它们对程序性能的负面影响也是不可小觑的。
*/
func BenchmarkConcurrentMap(b *testing.B) {
	keyType := reflect.TypeOf(int32(2))
	elemType := keyType
	cmap := NewConcurrentMap(keyType, elemType)
	var key, elem int32
	fmt.Printf("N=%d.\n", b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		seed := int32(i)
		key = seed
		elem = seed << 10
		b.StartTimer()
		cmap.Put(key, elem)
		_ = cmap.Get(key)
		b.StopTimer()
		b.SetBytes(8)
		b.StartTimer()
	}
	ml := cmap.Len()
	b.StopTimer()
	mapType := fmt.Sprintf("ConcurrentMap<%s, %s>",
		keyType.Kind().String(), elemType.Kind().String())
	b.Logf("The length of %s value is %d.\n", mapType, ml)
	b.StartTimer()
}

func testConcurrentMap(
	t *testing.T,
	newConcurrentMap func() ConcurrentMap,
	genKey func() interface{},
	genElem func() interface{},
	keyKind reflect.Kind,
	elemKind reflect.Kind) {
	mapType := fmt.Sprintf("ConcurrentMap<keyType=%s, elemType=%s>", keyKind, elemKind)
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			t.Errorf("Fatal Error: %s: %s\n", mapType, err)
		}
	}()
	t.Logf("Starting Test%s...", mapType)

	// Basic
	cmap := newConcurrentMap()
	expectedLen := 0
	if cmap.Len() != expectedLen {
		t.Errorf("ERROR: The length of %s value %d is not %d!\n",
			mapType, cmap.Len(), expectedLen)
		t.FailNow()
	}
	expectedLen = 5
	testMap := make(map[interface{}]interface{}, expectedLen)
	var invalidKey interface{}
	for i := 0; i < expectedLen; i++ {
		key := genKey()
		testMap[key] = genElem()
		if invalidKey == nil {
			invalidKey = key
		}
	}
	for key, elem := range testMap {
		oldElem, ok := cmap.Put(key, elem)
		if !ok {
			t.Errorf("ERROR: Put (%v, %v) to %s value %d is failing!\n",
				key, elem, mapType, cmap)
			t.FailNow()
		}
		if oldElem != nil {
			t.Errorf("ERROR: Already had a (%v, %v) in %s value %d!\n",
				key, elem, mapType, cmap)
			t.FailNow()
		}
		t.Logf("Put (%v, %v) to the %s value %v.",
			key, elem, mapType, cmap)
	}
	if cmap.Len() != expectedLen {
		t.Errorf("ERROR: The length of %s value %d is not %d!\n",
			mapType, cmap.Len(), expectedLen)
		t.FailNow()
	}
	for key, elem := range testMap {
		contains := cmap.Contains(key)
		if !contains {
			t.Errorf("ERROR: The %s value %v do not contains %v!",
				mapType, cmap, key)
			t.FailNow()
		}
		actualElem := cmap.Get(key)
		if actualElem == nil {
			t.Errorf("ERROR: The %s value %v do not contains %v!",
				mapType, cmap, key)
			t.FailNow()
		}
		t.Logf("The %s value %v contains key %v.", mapType, cmap, key)
		if actualElem != elem {
			t.Errorf("ERROR: The element of %s value %v with key %v do not equals %v!\n",
				mapType, actualElem, key, elem)
			t.FailNow()
		}
		t.Logf("The element of %s value %v to key %v is %v.",
			mapType, cmap, key, actualElem)
	}
	oldElem := cmap.Remove(invalidKey)
	if oldElem == nil {
		t.Errorf("ERROR: Remove %v from %s value %d is failing!\n",
			invalidKey, mapType, cmap)
		t.FailNow()
	}
	t.Logf("Removed (%v, %v) from the %s value %v.",
		invalidKey, oldElem, mapType, cmap)
	delete(testMap, invalidKey)

	// Type
	actualElemType := cmap.ElemType()
	if actualElemType == nil {
		t.Errorf("ERROR: The element type of %s value is nil!\n",
			mapType)
		t.FailNow()
	}
	actualElemKind := actualElemType.Kind()
	if actualElemKind != elemKind {
		t.Errorf("ERROR: The element type of %s value %s is not %s!\n",
			mapType, actualElemKind, elemKind)
		t.FailNow()
	}
	t.Logf("The element type of %s value %v is %s.",
		mapType, cmap, actualElemKind)
	actualKeyKind := cmap.KeyType().Kind()
	if actualKeyKind != elemKind {
		t.Errorf("ERROR: The key type of %s value %s is not %s!\n",
			mapType, actualKeyKind, keyKind)
		t.FailNow()
	}
	t.Logf("The key type of %s value %v is %s.",
		mapType, cmap, actualKeyKind)

	// Export
	keys := cmap.Keys()
	elems := cmap.Elems()
	pairs := cmap.ToMap()
	for key, elem := range testMap {
		var hasKey bool
		for _, k := range keys {
			if k == key {
				hasKey = true
			}
		}
		if !hasKey {
			t.Errorf("ERROR: The keys of %s value %v do not contains %v!\n",
				mapType, cmap, key)
			t.FailNow()
		}
		var hasElem bool
		for _, e := range elems {
			if e == elem {
				hasElem = true
			}
		}
		if !hasElem {
			t.Errorf("ERROR: The elems of %s value %v do not contains %v!\n",
				mapType, cmap, elem)
			t.FailNow()
		}
		var hasPair bool
		for k, e := range pairs {
			if k == key && e == elem {
				hasPair = true
			}
		}
		if !hasPair {
			t.Errorf("ERROR: The elems of %s value %v do not contains (%v, %v)!\n",
				mapType, cmap, key, elem)
			t.FailNow()
		}
	}

	// Clear
	cmap.Clear()
	if cmap.Len() != 0 {
		t.Errorf("ERROR: Clear %s value %d is failing!\n",
			mapType, cmap)
		t.FailNow()
	}
	t.Logf("The %s value %v has been cleared.", mapType, cmap)
}

func genRandString() string {
	var buff bytes.Buffer
	var prev string
	var curr string
	for i := 0; buff.Len() < 3; i++ {
		curr = string(genRandAZAscii())
		if curr == prev {
			continue
		}
		prev = curr
		buff.WriteString(curr)
	}
	return buff.String()
}

func genRandAZAscii() int {
	min := 65 // A
	max := 90 // Z
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}
