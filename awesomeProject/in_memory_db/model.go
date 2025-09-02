package in_memory_db

import (
	"errors"
	"time"
)

// Item represents a stored value
type Item struct {
	value     []byte
	expiresAt time.Time // zero means no expiry
	ver       uint64    // simple version for CAS
	size      int
}

// ErrKeyNotFound returned when key does not exist or expired
var ErrKeyNotFound = errors.New("key not found")
var ErrCASFailed = errors.New("cas failed")

// pair stored in LRU list
type kvPair struct {
	key  string
	item *Item
}

// Stats for DB
type Stats struct {
	Gets      uint64
	Sets      uint64
	Deletes   uint64
	Hits      uint64
	Misses    uint64
	Evictions uint64
	Bytes     uint64
}
