package ft

/**
使用testing测试包的测试文件名称需要以_test.go结尾，并且该测试文件需要与待测试的文件置于同一目录下。
功能测试函数名称必须以 Test为前缀，并且其后的字符串的第一个字符必须是大写，或者是数字。
唯一参数的类型必须是 *testing.T 类型的参数声明，如果没有按照此规则进行命名，则该函数在测试时不会被执行。

t.Log、t.Logf方法的作用，常规打印日志，测试通过 则不会打印
如果想查看常规测试日志 可以用go test -v
如果想让某个测试函数执行过程中 立即失败 可以调用t.FailNow方法
t.Fail() 表示测试失败，但是后续代码可执行
t.FailNow() 表示 当前函数立即停止，后续代码不再执行
想在测试失败的同时打印失败测试日志 直接调用t.Error方法或者t.Errorf方法
t.Fatal方法和t.Fatalf方法，它们的作用是在打印失败错误日志之后立即终止当前测试函数的执行并宣告测试失败。
更具体地说，这相当于它们在最后都调用了t.FailNow方法。
*/

/**
该函数负责测试User().Login()方法的功能，当其功能正常执行时返回nil，则调用t.Log()输出测试通过的信息，否则测试不通过，调用t.Error()输出错误信息。
*/
/*func TestUserLogin(t *testing.T) {
	err := User().Login("test1", "321")
	if err == nil {
		t.Log("Login test1 success")
	} else {
		t.Error(err)
	}
}*/
