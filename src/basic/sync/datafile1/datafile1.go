package datafile1

import (
	"errors"
	"io"
	"os"
	"sync"
)

/**
假设我们需要创建一个文件来存放数据。在同一时刻，可能会有多个Goroutine对这个文件进行写操作和读操作。每一次写操作都应该向这个文件写入若干个字节的数据。
这若干字节的数据应该作为独立的数据块存在。这就意味着，写操作之间不能彼此干扰，写入的内容之间也不能出现穿插和混淆的情况。另一方面，每一次读操作都应该从
这个文件中读取一个独立、完整的数据块。它们读取的数据块不能重复，且需要按顺序读取。例如，第一个读操作读取了数据块1，那么第二个操作就应该去读取数据块2，
以此类推。对于这些读操作是否可以同时进行，并不做要求。即使它们被同时进行，程序也应该分辨出它们的先后顺序。
*/

// 数据文件的接口类型
type DataFile interface {
	// 读取一个数据块
	Read() (rsn int64, d Data, err error)
	// 写入一个数据块
	Write(d Data) (wsn int64, err error)
	// 获取最后读取的数据块的序列号,这里所说的序列号相当于一个计数值，从1开始，得到当前已被读取的数据块的数量
	Rsn() int64
	// 获取最后写入的数据块的序列号,这里所说的序列号相当于一个计数值，从1开始，得到当前已被写入的数据块的数量
	Wsn() int64
	// 获取数据块的长度
	DataLen() uint32
}

// 数据文件的实现类型
type myDataFile struct {
	// 文件
	f *os.File
	// 被用于文件的读写锁。
	fMutex sync.RWMutex
	// 写操作需要用到的偏移量
	wOffset int64
	// 读操作需要用到的偏移量
	rOffset int64
	// 写操作用到的互斥锁
	wMutex sync.Mutex
	// 读操作需要用到的互斥锁
	rMutex sync.Mutex
	// 数据块长度
	dataLen uint32
}

type Data []byte

func NewDataFile(path string, dataLen uint32) (DataFile, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	if dataLen == 0 {
		return nil, errors.New("invalid data length")
	}
	// wOffset和rOffset的零值就是0，fMutex、wMutex和rMutex的零值就是可用的锁
	df := &myDataFile{
		f:       f,
		dataLen: dataLen,
	}
	return df, nil
}

func (df *myDataFile) Read() (rsn int64, d Data, err error) {
	// 读取并更新偏移量
	var offset int64
	// 这里使用互斥锁的目的是为了在多个Goroutine执行的情况下获取不重复且正确的读偏移量
	df.rMutex.Lock()
	offset = df.rOffset
	df.rOffset += int64(df.dataLen)
	df.rMutex.Unlock()
	// 读取一个数据块
	rsn = offset / int64(df.dataLen)
	bytes := make([]byte, df.dataLen)
	for {
		df.fMutex.RLock()
		_, err = df.f.ReadAt(bytes, offset)
		if err != nil {
			// 出现EOF的时候继续尝试获取同一个数据块，知道获取成功为止。这是为了避免在读Goroutine多于写Goroutine的情况下出现漏读的问题
			if err == io.EOF {
				// 必须每次循环到这都要解锁，否则写锁定操作将永远不会成功，且相应的Goroutine也会被一直阻塞
				df.fMutex.RUnlock()
				continue
			}
			df.fMutex.RUnlock()
			return
		}
		d = bytes
		df.fMutex.RUnlock()
		return
	}
}

func (df *myDataFile) Write(d Data) (wsn int64, err error) {
	// 读取并更新写偏移量
	var offset int64
	df.wMutex.Lock()
	offset = df.wOffset
	df.wOffset += int64(df.dataLen)
	df.wMutex.Unlock()

	// 写入一个数据块
	wsn = offset / int64(df.dataLen)
	var bytes []byte
	if len(d) > int(df.dataLen) {
		bytes = d[0:df.dataLen]
	} else {
		bytes = d
	}
	df.fMutex.Lock()
	defer df.fMutex.Unlock()
	_, err = df.f.Write(bytes)
	return
}

func (df *myDataFile) Rsn() int64 {
	df.rMutex.Lock()
	defer df.rMutex.Unlock()
	return df.rOffset / int64(df.dataLen)
}

func (df *myDataFile) Wsn() int64 {
	df.wMutex.Lock()
	defer df.wMutex.Unlock()
	return df.wOffset / int64(df.dataLen)
}

func (df *myDataFile) DataLen() uint32 {
	return df.dataLen
}
