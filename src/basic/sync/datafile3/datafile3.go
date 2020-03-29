package datafile3

import (
	"errors"
	"io"
	"os"
	"sync"
	"sync/atomic"
)

/**
通过用原子操作来代替锁的改造，程序性能会有一定的提升，因为原子操作是底层硬件支持，而锁操作是由操作系统提供的API实现。
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
	rCond *sync.Cond
	// 写操作需要用到的偏移量
	wOffset int64
	// 读操作需要用到的偏移量
	rOffset int64
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
	df.rCond = sync.NewCond(df.fMutex.RLocker())
	return df, nil
}

func (df *myDataFile) Read() (rsn int64, d Data, err error) {
	// 读取并更新偏移量
	var offset int64
	for {
		// 这样读取的原因是：在32位计算机上对64位整数进行操作的时候存在并发安全问题
		offset = atomic.LoadInt64(&df.rOffset)
		// 通过CAS更新
		if atomic.CompareAndSwapInt64(&df.rOffset, offset, offset+int64(df.dataLen)) {
			break
		}
	}
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
				df.rCond.Wait()
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
	for {
		// 这样读取的原因是：在32位计算机上对64位整数进行操作的时候存在并发安全问题
		offset = atomic.LoadInt64(&df.wOffset)
		// 通过CAS更新
		if atomic.CompareAndSwapInt64(&df.wOffset, offset, offset+int64(df.dataLen)) {
			break
		}
	}

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
	df.rCond.Signal()
	return
}

/**
这里读取需要加锁的原因是：在32位计算机上对64位整数进行操作的时候存在并发安全问题。
原因就是：线程切换带来的原子性问题。
64位的整数，32位机器读或写这个变量时得把人家咔嚓分成两个32位操作，可能一个线程读了某个值的高32位，低32位已经被另一个线程改了。
*/
func (df *myDataFile) Rsn() int64 {
	offset := atomic.LoadInt64(&df.rOffset)
	return offset / int64(df.dataLen)
}

func (df *myDataFile) Wsn() int64 {
	offset := atomic.LoadInt64(&df.wOffset)
	return offset / int64(df.dataLen)
}

func (df *myDataFile) DataLen() uint32 {
	return df.dataLen
}
