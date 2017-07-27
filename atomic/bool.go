// github.com/Vonng/gopher/atomic/bool.go provides an atomic bool type
package atomic

import "sync/atomic"

// AtomicBool use int32 atomic primitives to implement an atomic bool type
type AtomicBool struct{ flag int32 }

// Set to given bool value
func (ab *AtomicBool) Set(value bool) {
	var i int32 = 0
	if value {
		i = 1
	}
	atomic.StoreInt32(&(ab.flag), int32(i))
}

// Load bool value
func (ab *AtomicBool) Get() bool {
	if atomic.LoadInt32(&(ab.flag)) != 0 {
		return true
	}
	return false
}
