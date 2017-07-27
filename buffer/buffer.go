// buffers can be simply implemented by buffered channel. while this approach
// has it's limitations: Raw channel will raise runtime exception when:
// * sending value to closed channel
// * close a closed channel
// thus interface Buffer aims at transform runtime panic into normal errors
package buffer

import (
	"sync"
	"sync/atomic"
)

/**************************************************************
* interface: Buffer
**************************************************************/
// Buffer : a FIFO list
type Buffer interface {
	Cap() uint32
	Len() uint32
	// Put will put datum into buffer without block.
	// return non-nil error if buffer is already closed,
	Put(datum interface{}) (bool, error)
	// Get will get datum from buffer without block
	// return non-nil error if buffer is already closed,
	Get() (interface{}, error)
	// Close will close buffer
	// return false if buffer already closed, otherwise true.
	Close() bool
	// Closed indicate close status of buffer
	Closed() bool
}

/**************************************************************
* struct: defaultBuffer
**************************************************************/
// defaultBuffer : the default implementation of interface Buffer
type defaultBuffer struct {
	// ch : the low-level buffered channel
	ch chan interface{}
	// closed : bool-like closing status. 1 stand for true(closed)
	closed uint32
	// lock : eliminate race-condition on closing buffer
	lock sync.RWMutex
}

// [PUBLIC]
// NewBuffer will create a buffer with given size parameter
func NewBuffer(size uint32) (Buffer, error) {
	if size == 0 {
		return nil, ErrInvalidBufferSize
	}
	return &defaultBuffer{ch: make(chan interface{}, size) }, nil
}

// defaultBuffer_Cap proxy cap(chan)
func (buf *defaultBuffer) Cap() uint32 {
	return uint32(cap(buf.ch))
}

// defaultBuffer_Len proxy len(chan)
func (buf *defaultBuffer) Len() uint32 {
	return uint32(len(buf.ch))
}

// defaultBuffer_Put implements Buffer.Put. May race with Close
func (buf *defaultBuffer) Put(datum interface{}) (ok bool, err error) {
	buf.lock.RLock()
	defer buf.lock.RUnlock()
	if buf.Closed() {
		return false, ErrClosedBuffer
	}
	select {
	case buf.ch <- datum:
		ok = true
	default:
		ok = false
	}
	return
}

// defaultBuffer_Get will fetch a datum without block
func (buf *defaultBuffer) Get() (interface{}, error) {
	select {
	case datum, ok := <-buf.ch:
		if !ok {
			return nil, ErrClosedBuffer
		}
		return datum, nil
	default:
		return nil, nil
	}
}

// defaultBuffer_Close may race with Put
// despite returning flag, the buffer is ensure closed
func (buf *defaultBuffer) Close() bool {
	// CAS: if actually not closed(0), then set flag to closed(1)
	if atomic.CompareAndSwapUint32(&buf.closed, 0, 1) {
		buf.lock.Lock()
		close(buf.ch)
		buf.lock.Unlock()
		return true
	}
	// already closed(flag=1)
	return false
}

// defaultBuffer_Closed indicate whether buffer is closed
func (buf *defaultBuffer) Closed() bool {
	if atomic.LoadUint32(&buf.closed) == 0 {
		return false
	}
	return true
}
