package datafile2

import (
	"errors"
	"io"
	"os"
	"sync"
)

/**
条件变量
在Go语言中，sync.Cond代表了条件变量，它总是要与互斥量组合使用。类型*sync.Cond的方法集合中有3个方法：Wait方法、Signal方法和Broadcast方法。
它们分别代表了等待通知、单发通知和广播通知操作。可以和Java中的wait、notify和notifyAll对比起来看。
方法Wait会自动地对与该条件变量关联的那个锁进行解锁，并且使调用方所在的Goroutine被阻塞。一旦该方法收到通知，就会试图再次锁定该锁。如果锁定成功，它就
会唤醒那个被阻塞的Goroutine。否则，该方法会等待下一个通知，那个Goroutine也会继续被阻塞。而方法Signal和Broadcast的作用都是发送通知以唤醒正在为此
被阻塞的Goroutine。不同的是，前者的目标只有一个，而后者的目标是所有。
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
	// 读操作需要用到的条件变量
	rcond *sync.Cond
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
	df.rcond = sync.NewCond(df.fMutex.RLocker())
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
	df.fMutex.RLock()
	defer df.fMutex.RUnlock()
	for {
		_, err = df.f.ReadAt(bytes, offset)
		if err != nil {
			// 出现EOF的时候继续尝试获取同一个数据块，知道获取成功为止。这是为了避免在读Goroutine多于写Goroutine的情况下出现漏读的问题
			if err == io.EOF {
				// 等待知道写操作发送通知唤醒
				df.rcond.Wait()
				continue
			}
			return
		}
		d = bytes
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
	df.rcond.Signal()
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
