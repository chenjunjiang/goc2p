package main

import "fmt"

func getContent() string {
	return "Ruby"
}

func main() {
	/*content := "Go"
	switch content {
	default:
		fmt.Println("Unknown language")
	case "Python":
		fmt.Println("A interpreted language")
	case "Go":
		fmt.Println("A compiled language")
	}*/
	/*switch content := getContent(); content {
	default:
		fmt.Println("Unknown language")
	case "Python":
		fmt.Println("A interpreted language")
	case "Go":
		fmt.Println("A compiled language")
	}*/
	/*switch content := getContent(); content {
	default:
		fmt.Println("Unknown language")
	case "Ruby", "Python":
		fmt.Println("A interpreted language")
	case "C", "Java", "Go":
		fmt.Println("A compiled language")
	}*/
	/*switch content := getContent(); content {
	default:
		fmt.Println("Unknown language")
	case "Ruby":
		fallthrough // fallthrough语句会将流程控制权转移到下一条case语句上，它只能在case语句列表的最后一句。即使content是Ruby，也会输出A interpreted language
	case "Python":
		fmt.Println("A interpreted language")
	case "C", "Java", "Go":
		fmt.Println("A compiled language")
	}*/
	// 类型switch语句
	var v interface{} = 12
	/*switch v.(type) {
	case string: // 如果v的类型是
		fmt.Printf("The string is '%s'.\n", v.(string))
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64:
		fmt.Printf("The integer is %d.\n", v)
	default:
		fmt.Printf("Unsupported value.(type=%T)\n", v)
	}*/
	switch i := v.(type) {
	case string: // 如果v的类型是
		fmt.Printf("The string is '%s'.\n", i)
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64:
		fmt.Printf("The integer is %d.\n", i)
	default:
		fmt.Printf("Unsupported value.(type=%T)\n", i)
	}
	/*number := 99
	// 在switch表达式缺失的情况下，该switch语句的判定目标会被视为布尔类型值true
	switch {
	case number < 100:
		number++
	case number < 200:
		number--
	default:
		number -= 2
	}
	fmt.Println(number)*/
	switch number := 123; {
	case number < 100:
		number++
	case number < 200:
		number--
	default:
		number -= 2
	}
}
