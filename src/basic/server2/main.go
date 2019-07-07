package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

var mu sync.Mutex
var count int

/**
通过请求的URL来决定哪一个被调用，对于传入的请求，服务器在不同的goroutine中运行该处理函数，这样它可以同时处理处理多个请求。然而，如果两个并发的请求
试图同时更新计数值count，它可能会不一致地增加，程序会产生一个严重的竞态bug。为了避免该问题，必须确保最多只有一个goroutine在同一时间访问变量，这正是
mu.Lock()和mu.Unlock()语句的作用。
*/
func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/count", counter)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func counter(writer http.ResponseWriter, request *http.Request) {
	mu.Lock()
	fmt.Fprintf(writer, "Count %d\n", count)
	mu.Unlock()
}

func handler(writer http.ResponseWriter, request *http.Request) {
	mu.Lock()
	count++
	mu.Unlock()
	fmt.Fprintf(writer, "URL.Path = %q\n", request.URL.Path)
}
