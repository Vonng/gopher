// redis_pool is a pool in hbase
package buffer

type Serializable interface {
	Serialize() ([]byte, error)
	Deserialize(s []byte) (interface{}, error)
}

// RedisPool implement pool for serializable data structure
type RedisPool interface {
	// Total fetch current number of item in pool
	Total() uint64
	Exist(datum interface{}) error
	// Get will fetch an item from pool
	// return non-nil err if pool is already closed
	Get() (datum interface{}, err error)
	// Close will close the pool
	// return false if pool already closed, else true
	Close() bool
	// Closed indicate pool's closing status
	Closed() bool
}
