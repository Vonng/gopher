package atomic

// Semaphore
// https://en.wikipedia.org/wiki/Semaphore_(programming)
// A variable used to control access to a common resource
// by multiple routines in a concurrent system.
// Implemented using gol raw channel
type Semaphore chan struct{}

// There is not NewSemaphore. Instead, Just using builtin `make`
// make(Semaphore,N) will

// Semaphore_P acquire n resources
func (s Semaphore) P(n int) {
	e := struct{}{};
	for i := 0; i < n; i++ {
		s <- e
	}
}

// Semaphore_V release n resources
func (s Semaphore) V(n int) {
	for i := 0; i < n; i++ {
		<-s
	}
}

// Semaphore_Inc release(add) 1 resource
func (s Semaphore) Inc() {
	s <- struct{}{};
}

// Semaphore_Dec acquire(minus) 1 resource
func (s Semaphore) Dec() {
	<-s
}
