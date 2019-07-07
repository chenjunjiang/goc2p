package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

/**
这个进程使用的时间是不超过最耗时任务的时间，而不是所有任务总的时间。

chenjunjiang@chenjunjiang-B85-HD3:~/go_workspace/goc2p/src/basic/fetchall$ ./fetchall https://golang.org http://gopl.io https://godoc.org
1.34s    4154 http://gopl.io
1.60s    6797 https://godoc.org
Get https://golang.org: dial tcp 216.239.37.1:443: i/o timeout
30.00s elapsed

chenjunjiang@chenjunjiang-B85-HD3:~/go_workspace/goc2p/src/basic/fetchall$ ./fetchall https://golang.google.cn/ http://gopl.io https://godoc.org
0.53s    5530 https://golang.google.cn/
2.78s    4154 http://gopl.io
7.89s    6797 https://godoc.org
7.89s elapsed

 */
func main() {
	start := time.Now()
	ch := make(chan string)
	for _, url := range os.Args[1:] {
		go fetch(url, ch) // 启动一个goroutine
	}
	// 当一个goroutine试图在一个通道上进行发送或接收操作时，它会阻塞，直到另一个goroutine试图接收或发送操作才传递值，并开始处理两个goroutine。
	// 本例中，每一个fetch在通道ch上发送一个值，main函数接收它们。由main来处理所有的输出确保了每个goroutine作为一个整体单元处理，这样就避免了
	// 两个goroutine同时完成造成输出交织所带来的风险
	for range os.Args[1:] {
		fmt.Println(<-ch) // 从通道ch接收
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err) // 发送到通道ch
		return
	}

	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs %7d %s", secs, nbytes, url)
}
