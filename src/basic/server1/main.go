package main

import (
	"fmt"
	"log"
	"net/http"
)

/**
chenjunjiang@chenjunjiang-B85-HD3:~/go_workspace/goc2p/src/basic/fetch$ ./fetch  http://localhost:8000
URL.Path = "/"
chenjunjiang@chenjunjiang-B85-HD3:~/go_workspace/goc2p/src/basic/fetch$ ./fetch  http://localhost:8000/help
URL.Path = "/help"
 */
func main() {
	http.HandleFunc("/", handler) // 回声请求调用处理程序
	// 启动服务器监听8000端口
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

// 处理程序回显请求URL request的路径部分
func handler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "URL.Path = %q\n", request.URL.Path)
}
