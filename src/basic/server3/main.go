package main

import (
	"fmt"
	"log"
	"net/http"
)

/**
接收请求的消息头和表单数据
*/
func main() {
	http.HandleFunc("/", handler) // 回声请求调用处理程序
	// 启动服务器监听8000端口
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

// 处理程序回显请求http请求
func handler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "%s %s %s\n", request.Method, request.URL, request.Proto)
	for k, v := range request.Header {
		fmt.Fprintf(writer, "Header[%q] = %q\n", k, v)
	}
	fmt.Fprintf(writer, "Host = %q\n", request.Host)
	fmt.Fprintf(writer, "RemoteAddr = %q\n", request.RemoteAddr)
	if err := request.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range request.Form {
		fmt.Fprintf(writer, "Form[%q] = %q\n", k, v)
	}
}
