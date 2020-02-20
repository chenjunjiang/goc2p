package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

/**
管道是一种半双工（或者说是单向的）通讯方式，它只能被用于父进程与子进程以及同祖先的子进程之间的通讯。
Go语言是支持管道的，通过标准库代码包os/exec中的API，我们可以执行操作系统命令并在此之上建立管道。
*/

func main() {
	cmdo := exec.Command("echo", "-n", "My first command from golang.")
	// 创建一个获取此命令输出的管道
	stdouto, err := cmdo.StdoutPipe()
	if err != nil {
		fmt.Printf("Error: Can not obtain the stdout pipe for command No.0: %s\n", err)
		return
	}
	// 启动操作系统命令
	if err := cmdo.Start(); err != nil {
		fmt.Printf("Error: The command No.0 can not be startup: %s\n", err)
	}
	/*var outputBuf bytes.Buffer
	for {
		tempOutput := make([]byte, 5)
		n, err := stdouto.Read(tempOutput)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Printf("Error: Can not read data from the pipe: %s\n", err)
				return
			}
		}
		if n > 0 {
			outputBuf.Write(tempOutput[:n])
		}
	}
	fmt.Printf("%s\n", outputBuf.String())*/

	// 一开始就使用带缓冲的读取器从输出管道中读取数据
	// 这个缓冲器在默认情况下会携带一个长度为4096的缓冲区。缓冲区的长度代表我们一次可以读取的字节的最大数量
	outputBuf := bufio.NewReader(stdouto)
	// 第二个返回值表明当前行是否还未被读完,如果为false，我们依然可以利用for语句来读取剩余的数据
	outoutO, _, err := outputBuf.ReadLine()
	if err != nil {
		fmt.Printf("Error: Can not read data from the pipe: %s\n", err)
		return
	}
	fmt.Printf("%s\n", string(outoutO))

	fmt.Println("-----------------------------")

	// 管道是一个单向数据通道，它可以把一个命令的输出作为另一个命令的输入
	cmd1 := exec.Command("ps", "aux")
	cmd2 := exec.Command("grep", "java")
	// 在cmd1上建立一个输出管道
	stdout1, err := cmd1.StdoutPipe()
	if err != nil {
		fmt.Printf("Error: Can not obtain the stdout pipe for command No.0: %s\n", err)
		return
	}
	if err := cmd1.Start(); err != nil {
		fmt.Printf("Error: The command Can not be startup: %s\n", err)
		return
	}
	// 在cmd2上建立一个输入管道，并把与cmd1连接的输出管道中的数据全部写入到这个输入管道中
	outputBuf1 := bufio.NewReader(stdout1)
	stdin2, err := cmd2.StdinPipe()
	if err != nil {
		fmt.Printf("Error: Can not obtain the stdin pipe for command: %s\n", err)
		return
	}
	// 把outputBuf1中缓存的数据全部写入到Writer中，这就等于把第一个命令的输出内容通过管道传递给了第二个命令
	outputBuf1.WriteTo(stdin2)
	// 启动cmd2并关闭与它连接的输入管道，以完成数据的传递
	var outputBuf2 bytes.Buffer
	cmd2.Stdout = &outputBuf2
	if err := cmd2.Start(); err != nil {
		fmt.Printf("Error: The command Can not be startup: %s\n", err)
		return
	}
	err = stdin2.Close()
	if err != nil {
		fmt.Printf("Error: Can not close the stdio pipe: %s\n", err)
		return
	}
	// 为了获取cmd2的所有输出内容，我们需要等待它运行结束后，再去查看缓冲区outputBuf2
	if err := cmd2.Wait(); err != nil {
		fmt.Printf("Error: Can not wait for the command: %s\n", err)
	}
	// fmt.Printf("%s\n", outputBuf2.Bytes())

	// fileBasedPipe()
	inMemorySyncPipe()
}

func fileBasedPipe() {
	// 创建命名管道，命名管道默认会在其中一端还未就绪的时候阻塞另一端的进程
	reader, writer, err := os.Pipe()
	if err != nil {
		fmt.Printf("Error: Can not create the named pipe: %s\n", err)
	}
	go func() {
		output := make([]byte, 100)
		// 从管道读取数据
		n, err := reader.Read(output)
		if err != nil {
			fmt.Printf("Error: Can not read data from the named pipe: %s\n", err)
		}
		fmt.Printf("Read %d byte(s). [file-based pipe]\n", n)
		fmt.Println(string(output[:n])) // ABCDEFGHIJKLMNOPQRSTUVWXYZ
	}()
	input := make([]byte, 26)
	for i := 65; i <= 90; i++ {
		input[i-65] = byte(i)
	}
	// 向管道写入数据
	fmt.Println(string(input)) // ABCDEFGHIJKLMNOPQRSTUVWXYZ
	n, err := writer.Write(input)
	if err != nil {
		fmt.Printf("Error: Can not write data to the named pipe: %s\n", err)
	}
	fmt.Printf("Written %d byte(s). [file-based pipe]\n", n)
	time.Sleep(200 * time.Millisecond)
}

func inMemorySyncPipe() {
	/**
	命名管道可以被多路复用。所以， 当多个输入端同时写入数据的时候我们就不得不需要考虑操作原子性的问题。操作系统提供的管道是不提供原子操作支持的。
	为此，Go语言在标准库代码包io中提供了一个被存于内存中的、有原子性操作保证的管道，简称内存管道。
	*/
	// 创建内存管道
	reader, writer := io.Pipe()
	go func() {
		output := make([]byte, 100)
		n, err := reader.Read(output)
		if err != nil {
			fmt.Printf("Error: Can not read data from the named pipe: %s\n", err)
		}
		fmt.Printf("Read %d byte(s). [in-memory pipe]\n", n)
	}()
	input := make([]byte, 26)
	for i := 65; i <= 90; i++ {
		input[i-65] = byte(i)
	}
	n, err := writer.Write(input)
	if err != nil {
		fmt.Printf("Error: Can not write data to the named pipe: %s\n", err)
	}
	fmt.Printf("Written %d byte(s). [in-memory pipe]\n", n)
	time.Sleep(200 * time.Millisecond)
}
