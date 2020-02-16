package set

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func testSetLenAndContains(t *testing.T, newSet func() Set, typeName string) {
	t.Logf("Starting Test%sLenAndContains...", typeName)
	set, expectedElemMap := genRandSet(newSet)
	t.Logf("Got a %s value: %v.", typeName, set)
	expectedLen := len(expectedElemMap)
	if set.Len() != expectedLen {
		t.Errorf("ERROR: The length of %s value %d is not %d!\n",
			typeName, set.Len(), expectedLen)
		// 立即终止当前测试函数的执行,这会使得当前的测试程序转而去执行其它的测试函数
		t.FailNow()
	}
	t.Logf("The length of %s value is %d.\n", typeName, set.Len())
	for k := range expectedElemMap {
		if !set.Contains(k) {
			t.Errorf("ERROR: The %s value %v do not contains %v!",
				set, typeName, k)
			t.FailNow()
		}
	}
}

func testSetAdd(t *testing.T, newSet func() Set, typeName string) {
	t.Logf("Starting Test%sAdd...", typeName)
	set := newSet()
	var randElem interface{}
	var result bool
	expectedElemMap := make(map[interface{}]bool)
	for i := 0; i < 5; i++ {
		randElem = genRandElement()
		t.Logf("Add %v to the %s value %v.\n", randElem, typeName, set)
		result = set.Add(randElem)
		if expectedElemMap[randElem] && result {
			t.Errorf("ERROR: The element adding (%v => %v) is successful but should be failing!\n",
				randElem, set)
			t.FailNow()
		}
		if !expectedElemMap[randElem] && !result {
			t.Errorf("ERROR: The element adding (%v => %v) is failing!\n",
				randElem, set)
			t.FailNow()
		}
		expectedElemMap[randElem] = true
	}
	t.Logf("The %s value: %v.", typeName, set)
	expectedLen := len(expectedElemMap)
	if set.Len() != expectedLen {
		t.Errorf("ERROR: The length of %s value %d is not %d!\n",
			typeName, set.Len(), expectedLen)
		t.FailNow()
	}
	t.Logf("The length of %s value is %d.\n", typeName, set.Len())
	for k := range expectedElemMap {
		if !set.Contains(k) {
			t.Errorf("ERROR: The %s value %v do not contains %v!",
				typeName, set, k)
			t.FailNow()
		}
	}
}

func testSetRemove(t *testing.T, newSet func() Set, typeName string) {
	t.Logf("Starting Test %s Remove...", typeName)
	set, expectedElemMap := genRandSet(newSet)
	t.Logf("Got a %s value: %v.", typeName, set)
	t.Logf("The length of %s value is %d.\n", typeName, set.Len())
	var number int
	for k, _ := range expectedElemMap {
		if number%2 == 0 {
			t.Logf("Remove %v from the HashSet value %v.\n", k, set)
			set.Remove(k)
			if set.Contains(k) {
				t.Errorf("ERROR: The element removing (%v => %v) is failing!\n",
					k, set)
				t.FailNow()
			}
			delete(expectedElemMap, k)
		}
		number++
	}
	expectedLen := len(expectedElemMap)
	if set.Len() != expectedLen {
		t.Errorf("ERROR: The length of %v value %d is not %d!\n", typeName, set.Len(), expectedLen)
		t.FailNow()
	}
	t.Logf("The length of %s value is %d.\n", typeName, set.Len())
	// 删除元素之后剩下的元素是否是期望的
	for _, v := range set.Elements() {
		if !expectedElemMap[v] {
			t.Errorf("ERROR: The HashSet value %v contains %v but should not contains!", set, v)
			t.FailNow()
		}
	}
}

func testSetClear(t *testing.T, newSet func() Set, typeName string) {
	t.Logf("Starting Test%sClear...", typeName)
	set, _ := genRandSet(newSet)
	t.Logf("Got a %s value: %v.", typeName, set)
	t.Logf("The length of %s value is %d.\n", typeName, set.Len())
	t.Logf("Clear the HashSet value %v.\n", set)
	set.Clear()
	expectedLen := 0
	if set.Len() != expectedLen {
		t.Errorf("ERROR: The length of %v value %d is not %d!\n", typeName, set.Len(), expectedLen)
		t.FailNow()
	}
	t.Logf("The length of %s value is %d.\n", typeName, set.Len())
}

func testSetElements(t *testing.T, newSet func() Set, typeName string) {
	t.Logf("Starting Test %s Elements...", typeName)
	set, expectedElemMap := genRandSet(newSet)
	t.Logf("Got a %s value: %v.", typeName, set)
	t.Logf("The length of %s value is %d.\n", typeName, set.Len())
	elems := set.Elements()
	t.Logf("The elements of %s value is %v.\n", typeName, elems)
	expectedLen := len(expectedElemMap)
	if len(elems) != expectedLen {
		t.Errorf("ERROR: The length of HashSet value %d is not %d!\n", len(elems), expectedLen)
		t.FailNow()
	}
	t.Logf("The length of elements is %d.\n", len(elems))
	for _, v := range elems {
		if !expectedElemMap[v] {
			t.Errorf("ERROR: The elements %v contains %v but should not contains!", set, v)
			t.FailNow()
		}
	}
}

func testSetString(t *testing.T, newSet func() Set, typeName string) {
	t.Logf("Starting Test %s String...", typeName)
	set, _ := genRandSet(newSet)
	t.Logf("Got a %s value: %v.", typeName, set)
	setStr := set.String()
	t.Logf("The string of %s value is %s.\n", typeName, setStr)
	var elemStr string
	for _, v := range set.Elements() {
		elemStr = fmt.Sprintf("%v", v)
		if !strings.Contains(setStr, elemStr) {
			t.Errorf("ERROR: The string of %s value %s do not contains %s!",
				typeName, setStr, elemStr)
			t.FailNow()
		}
	}
}

/**
生成随机的测试对象
*/
func genRandSet(newSet func() Set) (set Set, elemMap map[interface{}]bool) {
	set = newSet()
	elemMap = make(map[interface{}]bool)
	var enough bool
	for !enough {
		e := genRandElement()
		set.Add(e)
		elemMap[e] = true
		if len(elemMap) >= 3 {
			enough = true
		}
	}
	return
}

func genRandElement() interface{} {
	seed := rand.Int63n(5)
	switch seed {
	case 0:
		return genRandInt()
	case 1:
		return genRandString()
	case 2:
		return struct {
			num int64
			str string
		}{genRandInt(), genRandString()}
	default:
		const length = 2
		arr := new([length]interface{})
		for i := 0; i < length; i++ {
			if i%2 == 0 {
				arr[i] = genRandInt()
			} else {
				arr[i] = genRandString()
			}
		}
		return *arr
	}
}

func genRandString() string {
	var buff bytes.Buffer
	var prev string
	var curr string
	for buff.Len() < 3 {
		curr = string(genRandAZAscii())
		if curr == prev {
			continue
		} else {
			prev = curr
		}
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

func genRandInt() int64 {
	return rand.Int63n(10000)
}
