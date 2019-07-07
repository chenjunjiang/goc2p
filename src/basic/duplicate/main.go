package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	/*counts := make(map[string]int)
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		counts[input.Text()]++
	}

	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s", n, line)
		}
	}*/

	/*counts := make(map[string]int)
	input := bufio.NewScanner(os.Stdin)
	// for循环是怎么退出的呢？ 在linux上， 按ctrl + d吧， 让input.Scan函数返回false.
	for input.Scan() {
		s := input.Text()
		fmt.Printf("cur line is %s\n", s)
		counts[s]++
	}

	for line, n := range counts {
		// 大于1表示有重复
		if n > 1 {
			fmt.Printf("%d  %s\n", n, line);
		}
	}*/

	// 从文件中读取数据
	/**
	chenjunjiang@chenjunjiang-B85-HD3:~/go_workspace/goc2p/src/basic/duplicate$ go build
    chenjunjiang@chenjunjiang-B85-HD3:~/go_workspace/goc2p/src/basic/duplicate$ ls
    duplicate  main.go
    chenjunjiang@chenjunjiang-B85-HD3:~/go_workspace/goc2p/src/basic/duplicate$ ./duplicate /home/chenjunjiang/duplicate.txt
    2	qwert
	 */
	/*counts := make(map[string]int)
	files := os.Args[1:]
	if len(files) == 0 {
		countLines(os.Stdin, counts)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup: %v\n", err)
				continue
			}
			countLines(f, counts)
			f.Close()
		}
	}

	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}*/

	/**
	一次读取整个文件内容到内存中，一次性地分割所有行
	chenjunjiang@chenjunjiang-B85-HD3:~/go_workspace/goc2p/src/basic/duplicate$ go build
	chenjunjiang@chenjunjiang-B85-HD3:~/go_workspace/goc2p/src/basic/duplicate$ ./duplicate /home/chenjunjiang/duplicate.txt
	找到重复行: 2	qwert
	 */
	counts := make(map[string]int)
	for _, filename := range os.Args[1:] {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "dup3: %v\n", err)
			continue
		}
		for _, line := range strings.Split(string(data), "\n") {
			counts[line]++
		}
	}

	for line, n := range counts {
		if n > 1 {
			fmt.Printf("找到重复行: %d\t%s\n", n, line)
		}
	}
}

/*func countLines(f *os.File, counts map[string]int) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		counts[input.Text()]++
	}
}*/
