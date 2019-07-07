package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

/**
chenjunjiang@chenjunjiang-B85-HD3:~/go_workspace/goc2p/src/basic/fetch$ ./fetch  http://gopl.io
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN"
	  "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
  <meta name="go-import" content="gopl.io git https://github.com/adonovan/gopl.io"></meta>
  <title>The Go Programming Language</title>
  <script>
......

chenjunjiang@chenjunjiang-B85-HD3:~/go_workspace/goc2p/src/basic/fetch$ ./fetch  http://gopxx.io
fetch: Get http://gopxx.io: dial tcp: lookup gopxx.io: no such host
 */
func main() {
	for _, url := range os.Args[1:] {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			os.Exit(1)
		}
		b, err := ioutil.ReadAll(resp.Body)
		// 避免资源泄露
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
		}
		fmt.Printf("%s", b)
	}
}
