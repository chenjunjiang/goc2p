package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

func main() {
	go func() {
		time.Sleep(5 * time.Second)
		sigSendingDemo()
	}()
	signalHandleDemo()
}

func signalHandleDemo() {
	sigRecv1 := make(chan os.Signal, 1)
	// 类Unix操作系统下有两种信号既不能被自行处理也不会被忽略，SIGKILL和SIGSTOP，对它们的响应只能是执行系统默认操作
	sigs1 := []os.Signal{syscall.SIGINT, syscall.SIGQUIT}
	fmt.Printf("Set notification for %s...[sigRecv1]\n", sigs1)
	// 在接受到我们希望自行处理的信号之后，signal处理程序会把它封装成syscall.Signal类型的值并放到signal接收通道中。如果没有第二个参数，将会接收所有信号
	signal.Notify(sigRecv1, sigs1...)
	sigRecv2 := make(chan os.Signal, 1)
	sigs2 := []os.Signal{syscall.SIGQUIT}
	fmt.Printf("Set notification for %s...[sigRecv2]\n", sigs2)
	signal.Notify(sigRecv2, sigs2...)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		/**
		只要sigRecv1中存在元素值，for语句就会把它们按顺序地接收并赋值给变量sig，否则for语句就会被阻塞，并等待新的元素值被发送到sigRecv1中，
		在sigRecv1代表的通道类型值被关闭之后，for语句就会立即退出
		*/
		for sig := range sigRecv1 {
			fmt.Printf("Received a signal from sigRecv1: %s\n", sig)
		}
		fmt.Printf("End. [sigRecv1]\n")
		wg.Done()
	}()
	go func() {
		for sig := range sigRecv2 {
			fmt.Printf("Received a signal from sigRecv2: %s\n", sig)
		}
		fmt.Printf("End. [sigRecv2]\n")
		wg.Done()
	}()

	fmt.Println("Wait for 20 seconds... ")
	time.Sleep(20 * time.Second)
	fmt.Printf("Stop notification...")
	signal.Stop(sigRecv1)
	close(sigRecv1)
	fmt.Printf("done. [sigRecv1]\n")
	wg.Wait()
}

/**
发送信号
*/
func sigSendingDemo() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Fatal Error: %s\n", err)
			debug.PrintStack()
		}
	}()
	// ps aux | grep "mysignal" | grep -v "grep" | grep -v "go run" | awk '{print $2}'
	cmds := []*exec.Cmd{
		exec.Command("ps", "aux"),
		exec.Command("grep", "mysignal"),
		exec.Command("grep", "-v", "grep"),
		exec.Command("grep", "-v", "go run"),
		exec.Command("awk", "{print $2}"),
	}
	output, err := runCmds(cmds)
	if err != nil {
		fmt.Printf("Command Execution Error: %s\n", err)
		return
	}
	pids, err := getPids(output)
	if err != nil {
		fmt.Printf("PID Parsing Error: %s\n", err)
		return
	}
	fmt.Printf("Target PID(s):\n%v\n", pids)
	for _, pid := range pids {
		proc, err := os.FindProcess(pid)
		if err != nil {
			fmt.Printf("Process Finding Error: %s\n", err)
			return
		}
		sig := syscall.SIGQUIT
		fmt.Printf("Send signal '%s' to the process (pid=%d)...\n", sig, pid)
		err = proc.Signal(sig)
		if err != nil {
			fmt.Printf("Signal Sending Error: %s\n", err)
			return
		}
	}
}

func getPids(strs []string) ([]int, error) {
	pids := make([]int, 0)
	for _, str := range strs {
		pid, err := strconv.Atoi(strings.TrimSpace(str))
		if err != nil {
			return nil, err
		}
		pids = append(pids, pid)
	}
	return pids, nil
}

/**
第一命令的输出作为第二个命令的输入，第二个命令的输出又作为第三个命令的输入，以此类推
*/
func runCmds(cmds []*exec.Cmd) ([]string, error) {
	if cmds == nil || len(cmds) == 0 {
		return nil, errors.New("the cmd slice is invalid")
	}
	first := true
	var output []byte
	var err error
	for _, cmd := range cmds {
		fmt.Printf("Run command: %v ...\n", getCmdPlaintext(cmd))
		// 把前一个命令输出作为当前命令的输入
		if !first {
			var stdinBuf bytes.Buffer
			// 把output写入Buffer
			stdinBuf.Write(output)
			cmd.Stdin = &stdinBuf
		}
		var stdoutBuf bytes.Buffer
		cmd.Stdout = &stdoutBuf
		if err = cmd.Start(); err != nil {
			return nil, getError(err, cmd)
		}
		if err = cmd.Wait(); err != nil {
			return nil, getError(err, cmd)
		}
		output = stdoutBuf.Bytes()
		if !first {
			fmt.Printf("Output:\n%s\n", string(output))
		}
		if first {
			first = false
		}
	}
	lines := make([]string, 0)
	var outputBuf bytes.Buffer
	outputBuf.Write(output)
	for {
		// 读取数据直到遇到换行
		line, err := outputBuf.ReadBytes('\n')
		if err != nil {
			// 所有数据已读完
			if err == io.EOF {
				break
			} else {
				return nil, getError(err, nil)
			}
		}
		lines = append(lines, string(line))
	}
	return lines, nil
}

func getError(err error, cmd *exec.Cmd, extraInfo ...string) error {
	var errMsg string
	if cmd != nil {
		errMsg = fmt.Sprintf("%s  [%s %v]", err, (*cmd).Path, (*cmd).Args)
	} else {
		errMsg = fmt.Sprintf("%s", err)
	}
	if len(extraInfo) > 0 {
		errMsg = fmt.Sprintf("%s (%v)", errMsg, extraInfo)
	}
	return errors.New(errMsg)
}

func getCmdPlaintext(cmd *exec.Cmd) string {
	var buf bytes.Buffer
	buf.WriteString(cmd.Path)
	for _, arg := range cmd.Args[1:] {
		buf.WriteRune(' ')
		buf.WriteString(arg)
	}
	return buf.String()
}

/**
只测试signalHandleDemo，通过命令执行(当然也可以直接用idea执行)，然后分别键入Ctrl-c(对应SIGINT信号)和Ctrl-\(对应SIGQUIT信号)，在关闭sigRecv1后，只有sigRecv2能接收
SIGQUIT信号了，最后键入Ctrl-c，由于没有能处理SIGINT信号的通道了，所以当前进程直接被停止
chenjunjiang@chenjunjiang-B85-HD3:~/go_workspace/goc2p/src/basic/concurrency/signal$ go run mysignal.go
Set notification for [interrupt quit]...[sigRecv1]
Set notification for [quit]...[sigRecv2]
Wait for 2 seconds...
^CReceived a signal from sigRecv1: interrupt
^\Received a signal from sigRecv1: quit
Received a signal from sigRecv2: quit
Stop notification...done. [sigRecv1]
End. [sigRecv1]
^\Received a signal from sigRecv2: quit
^\Received a signal from sigRecv2: quit
^\Received a signal from sigRecv2: quit
^Csignal: interrupt
*/

/**
测试sigSendingDemo和signalHandleDemo，需要通过命令执行，通过idea执行不能找到mysignal进程， 只能找到/tmp/___go_build_basic_concurrency_signal。
通过go run命令执行，还会生成go run mysignal.go进程，所以通过grep -v go run过滤掉；如果不想过滤，那么就先通过go build生成可执行文件mysignal，
然后在执行这个可执行文件。
chenjunjiang@chenjunjiang-B85-HD3:~/go_workspace/goc2p/src/basic/concurrency/signal$ go run mysignal.go
Set notification for [interrupt quit]...[sigRecv1]
Set notification for [quit]...[sigRecv2]
Wait for 20 seconds...
Run command: /bin/ps aux ...
Run command: /bin/grep mysignal ...
Output:
chenjun+  1369  1.8  0.1 932164 17168 pts/0    Sl+  14:42   0:00 go run mysignal.go
chenjun+  1450  0.0  0.0 103532  1592 pts/0    Sl+  14:42   0:00 /tmp/go-build999450432/b001/exe/mysignal

Run command: /bin/grep -v grep ...
Output:
chenjun+  1369  1.8  0.1 932164 17168 pts/0    Sl+  14:42   0:00 go run mysignal.go
chenjun+  1450  0.0  0.0 103532  1592 pts/0    Sl+  14:42   0:00 /tmp/go-build999450432/b001/exe/mysignal

Run command: /bin/grep -v go run ...
Output:
chenjun+  1450  0.0  0.0 103532  1592 pts/0    Sl+  14:42   0:00 /tmp/go-build999450432/b001/exe/mysignal

Run command: /usr/bin/awk {print $2} ...
Output:
1450

Target PID(s):
[1450]
Send signal 'quit' to the process (pid=1450)...
Received a signal from sigRecv1: quit
Received a signal from sigRecv2: quit
Stop notification...done. [sigRecv1]
End. [sigRecv1]
^Csignal: interrupt
*/
