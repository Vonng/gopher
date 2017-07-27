// A buffer pool implemented inspired by leaky_buffer
// https://golang.org/doc/effective_go.html#leaky_buffer
// Buffer is good for common use, while:
// * Get/Put is non-block, and will return err when buffer is filled
// * It's size is fixed
// Pool is another wrapper for buffer, which makes Get/Put a blocking
// operation. And can auto size-up at runtime.
package buffer

import (
	"sync"
	"sync/atomic"
)

/**************************************************************
* interface: Pool
**************************************************************/

// Pool is in-memory pool of buffer, inspired by leaky buffer
type Pool interface {
	// BufferCap is standard buffer size of Buffer in Pool
	BufferCap() uint32
	// MaxBufferNumber indicate how many buffer can Pool have
	MaxBufferNumber() uint32
	// BufferNumber will get current buffer numbers
	BufferNumber() uint32
	// Total fetch current number of item in pool
	Total() uint64
	// Put will blocking put item into pool
	// returns non-nil err if pool is already closed
	Put(datum interface{}) error
	// Get will fetch an item from pool
	// return non-nil err if pool is already closed
	Get() (datum interface{}, err error)
	// Close will close the pool
	// return false if pool already closed, else true
	Close() bool
	// Closed indicate pool's closing status
	Closed() bool
}

/**************************************************************
* struct: defaultPool
**************************************************************/

// defaultPool is the default implementation of Pool 
type defaultPool struct {
	// bufferCap : store standard size of buffer
	bufferCap uint32
	// maxBufferNumber max allowed number of buffer in pool
	maxBufferNumber uint32
	// bufferNumber actual number of buffer in pool
	bufferNumber uint32
	// total : total items in pool  
	total uint64
	// bufCh : channel of Buffer, Buffer of Buffer
	bufCh chan Buffer
	// closed : Pool close status. 1 stand for true(closed)
	closed uint32
	lock   sync.RWMutex
}

// NewPool create a new Buffer Pool with given params
// bufferCap代表池内缓冲器的统一容量。
// 参数maxBufferNumber代表池中最多包含的缓冲器的数量。
func NewPool(bufferCap uint32, maxBufferNumber uint32) (Pool, error) {
	if bufferCap == 0 || maxBufferNumber == 0 {
		return nil, ErrInvalidPoolSize
	}

	bufCh := make(chan Buffer, maxBufferNumber)
	buf, _ := NewBuffer(bufferCap)
	bufCh <- buf
	return &defaultPool{
		bufferCap:       bufferCap,
		maxBufferNumber: maxBufferNumber,
		bufferNumber:    1,
		bufCh:           bufCh,
	}, nil
}

func (pool *defaultPool) BufferCap() uint32 {
	return pool.bufferCap
}

func (pool *defaultPool) MaxBufferNumber() uint32 {
	return pool.maxBufferNumber
}

func (pool *defaultPool) BufferNumber() uint32 {
	return atomic.LoadUint32(&pool.bufferNumber)
}

func (pool *defaultPool) Total() uint64 {
	return atomic.LoadUint64(&pool.total)
}

func (pool *defaultPool) Put(datum interface{}) (err error) {
	if pool.Closed() {
		return ErrClosedPool
	}
	var count uint32
	maxCount := pool.BufferNumber() * 5
	var ok bool
	for buf := range pool.bufCh {
		ok, err = pool.putData(buf, datum, &count, maxCount)
		if ok || err != nil {
			break
		}
	}
	return
}

// putData 用于向给定的缓冲器放入数据，并在必要时把缓冲器归还给池。
func (pool *defaultPool) putData(
	buf Buffer, datum interface{}, count *uint32, maxCount uint32) (ok bool, err error) {
	if pool.Closed() {
		return false, ErrClosedPool
	}
	defer func() {
		pool.lock.RLock()
		if pool.Closed() {
			atomic.AddUint32(&pool.bufferNumber, ^uint32(0))
			err = ErrClosedPool
		} else {
			pool.bufCh <- buf
		}
		pool.lock.RUnlock()
	}()
	ok, err = buf.Put(datum)
	if ok {
		atomic.AddUint64(&pool.total, 1)
		return
	}
	if err != nil {
		return
	}
	// 若因缓冲器已满而未放入数据就递增计数。
	(*count)++
	// 如果尝试向缓冲器放入数据的失败次数达到阈值，
	// 并且池中缓冲器的数量未达到最大值，
	// 那么就尝试创建一个新的缓冲器，先放入数据再把它放入池。
	if *count >= maxCount &&
		pool.BufferNumber() < pool.MaxBufferNumber() {
		pool.lock.Lock()
		if pool.BufferNumber() < pool.MaxBufferNumber() {
			if pool.Closed() {
				pool.lock.Unlock()
				return
			}
			newBuf, _ := NewBuffer(pool.bufferCap)
			newBuf.Put(datum)
			pool.bufCh <- newBuf
			atomic.AddUint32(&pool.bufferNumber, 1)
			atomic.AddUint64(&pool.total, 1)
			ok = true
		}
		pool.lock.Unlock()
		*count = 0
	}
	return
}

func (pool *defaultPool) Get() (datum interface{}, err error) {
	if pool.Closed() {
		return nil, ErrClosedPool
	}
	var count uint32
	maxCount := pool.BufferNumber() * 10
	for buf := range pool.bufCh {
		datum, err = pool.getData(buf, &count, maxCount)
		if datum != nil || err != nil {
			break
		}
	}
	return
}

// getData 用于从给定的缓冲器获取数据，并在必要时把缓冲器归还给池。
func (pool *defaultPool) getData(
	buf Buffer, count *uint32, maxCount uint32) (datum interface{}, err error) {
	if pool.Closed() {
		return nil, ErrClosedPool
	}
	defer func() {
		// 如果尝试从缓冲器获取数据的失败次数达到阈值，
		// 同时当前缓冲器已空且池中缓冲器的数量大于1，
		// 那么就直接关掉当前缓冲器，并不归还给池。
		if *count >= maxCount &&
			buf.Len() == 0 &&
			pool.BufferNumber() > 1 {
			buf.Close()
			atomic.AddUint32(&pool.bufferNumber, ^uint32(0))
			*count = 0
			return
		}
		pool.lock.RLock()
		if pool.Closed() {
			atomic.AddUint32(&pool.bufferNumber, ^uint32(0))
			err = ErrClosedPool
		} else {
			pool.bufCh <- buf
		}
		pool.lock.RUnlock()
	}()
	datum, err = buf.Get()
	if datum != nil {
		atomic.AddUint64(&pool.total, ^uint64(0))
		return
	}
	if err != nil {
		return
	}
	// 若因缓冲器已空未取出数据就递增计数。
	(*count)++
	return
}

func (pool *defaultPool) Close() bool {
	if !atomic.CompareAndSwapUint32(&pool.closed, 0, 1) {
		return false
	}
	pool.lock.Lock()
	defer pool.lock.Unlock()
	close(pool.bufCh)
	for buf := range pool.bufCh {
		buf.Close()
	}
	return true
}

func (pool *defaultPool) Closed() bool {
	if atomic.LoadUint32(&pool.closed) == 1 {
		return true
	}
	return false
}
