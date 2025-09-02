package in_memory_db

import (
	"container/list"
	"sync"
	"time"
)

// Simple, production-minded in-memory key-value database in Go.
// Features implemented:
// - Thread-safe Get/Set/Delete
// - TTL (key expiration) with background janitor
// - Capacity limit with LRU eviction
// - Atomic Compare-And-Set (CAS)
// - Simple transaction support (batch of ops applied atomically)
// - Basic Scan/Keys
// - Stats

// Usage: run `go run in_memory_db.go` to see a usage example in main.

// DB is the in-memory database
type DB struct {
	mu        sync.RWMutex
	data      map[string]*list.Element // map to list element for LRU
	lru       *list.List               // front = most recent, back = least
	capacity  int                      // max number of items (0 = unlimited)
	janitorCh chan struct{}
	closed    chan struct{}
	stats     Stats
}

// NewDB creates a new DB with optional capacity for LRU eviction and janitor interval for TTL cleanup
func NewDB(capacity int, janitorInterval time.Duration) *DB {
	db := &DB{
		data:      make(map[string]*list.Element),
		lru:       list.New(),
		capacity:  capacity,
		janitorCh: make(chan struct{}),
		closed:    make(chan struct{}),
	}

	if janitorInterval > 0 {
		go db.janitor(janitorInterval)
	}

	return db
}

// Set stores a value (replaces existing). ttlSeconds==0 means no expiry.
func (db *DB) Set(key string, value []byte, ttlSeconds int) {
	db.mu.Lock()
	defer db.mu.Unlock()

	var expires time.Time
	if ttlSeconds > 0 {
		expires = time.Now().Add(time.Duration(ttlSeconds) * time.Second)
	}

	if el, ok := db.data[key]; ok {
		// update existing
		kp := el.Value.(*kvPair)
		oldSize := kp.item.size
		kp.item.value = append([]byte(nil), value...)
		kp.item.expiresAt = expires
		kp.item.ver++
		kp.item.size = len(value)
		db.lru.MoveToFront(el)
		db.stats.Sets++
		db.stats.Bytes += uint64(kp.item.size - oldSize)
		return
	}

	item := &Item{value: append([]byte(nil), value...), expiresAt: expires, ver: 1, size: len(value)}
	el := db.lru.PushFront(&kvPair{key: key, item: item})
	db.data[key] = el
	db.stats.Sets++
	db.stats.Bytes += uint64(item.size)

	if db.capacity > 0 && db.lru.Len() > db.capacity {
		db.evictLRU()
	}
}

// Get fetches value by key, returns ErrKeyNotFound if not found or expired
func (db *DB) Get(key string) ([]byte, error) {
	db.mu.RLock() // need write lock to move element to front
	defer db.mu.RUnlock()

	db.stats.Gets++
	el, ok := db.data[key]
	if !ok {
		db.stats.Misses++
		return nil, ErrKeyNotFound
	}

	kp := el.Value.(*kvPair)
	if kp.item.expiresAt.IsZero() == false && time.Now().After(kp.item.expiresAt) {
		// expired: delete
		db.removeElement(el)
		db.stats.Misses++
		return nil, ErrKeyNotFound
	}

	// hit
	db.lru.MoveToFront(el)
	db.stats.Hits++
	return append([]byte(nil), kp.item.value...), nil
}

// Delete removes key
func (db *DB) Delete(key string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	el, ok := db.data[key]
	if !ok {
		return ErrKeyNotFound
	}
	db.stats.Deletes++
	db.removeElement(el)
	return nil
}

// CAS does compare-and-set based on version. If expectedVer==0 it acts like Set-if-not-exist
func (db *DB) CAS(key string, expectedVer uint64, newValue []byte, ttlSeconds int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	el, ok := db.data[key]
	if !ok {
		if expectedVer != 0 {
			return ErrCASFailed
		}
		// create
		var expires time.Time
		if ttlSeconds > 0 {
			expires = time.Now().Add(time.Duration(ttlSeconds) * time.Second)
		}
		item := &Item{value: append([]byte(nil), newValue...), expiresAt: expires, ver: 1, size: len(newValue)}
		el = db.lru.PushFront(&kvPair{key: key, item: item})
		db.data[key] = el
		db.stats.Sets++
		db.stats.Bytes += uint64(item.size)
		if db.capacity > 0 && db.lru.Len() > db.capacity {
			db.evictLRU()
		}
		return nil
	}

	kp := el.Value.(*kvPair)
	if kp.item.expiresAt.IsZero() == false && time.Now().After(kp.item.expiresAt) {
		// expired
		db.removeElement(el)
		if expectedVer != 0 {
			return ErrCASFailed
		}
		// create new
		var expires time.Time
		if ttlSeconds > 0 {
			expires = time.Now().Add(time.Duration(ttlSeconds) * time.Second)
		}
		item := &Item{value: append([]byte(nil), newValue...), expiresAt: expires, ver: 1, size: len(newValue)}
		el = db.lru.PushFront(&kvPair{key: key, item: item})
		db.data[key] = el
		db.stats.Sets++
		db.stats.Bytes += uint64(item.size)
		return nil
	}

	if kp.item.ver != expectedVer {
		return ErrCASFailed
	}

	oldSize := kp.item.size
	kp.item.value = append([]byte(nil), newValue...)
	kp.item.ver++
	kp.item.size = len(newValue)
	kp.item.expiresAt = time.Time{}
	if ttlSeconds > 0 {
		kp.item.expiresAt = time.Now().Add(time.Duration(ttlSeconds) * time.Second)
	}
	db.lru.MoveToFront(el)
	db.stats.Bytes += uint64(kp.item.size - oldSize)
	return nil
}

// Commit applies transaction atomically
func (db *DB) Commit(tx *Tx) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// apply all ops
	for _, op := range tx.oplist {
		switch op.Type {
		case OpSet:
			var expires time.Time
			if op.TTLSeconds > 0 {
				expires = time.Now().Add(time.Duration(op.TTLSeconds) * time.Second)
			}
			if el, ok := db.data[op.Key]; ok {
				kp := el.Value.(*kvPair)
				oldSize := kp.item.size
				kp.item.value = append([]byte(nil), op.Value...)
				kp.item.expiresAt = expires
				kp.item.ver++
				kp.item.size = len(op.Value)
				db.lru.MoveToFront(el)
				db.stats.Bytes += uint64(kp.item.size - oldSize)
				db.stats.Sets++
			} else {
				item := &Item{value: append([]byte(nil), op.Value...), expiresAt: expires, ver: 1, size: len(op.Value)}
				el := db.lru.PushFront(&kvPair{key: op.Key, item: item})
				db.data[op.Key] = el
				db.stats.Sets++
				db.stats.Bytes += uint64(item.size)
			}
			if db.capacity > 0 && db.lru.Len() > db.capacity {
				db.evictLRU()
			}
		case OpDelete:
			if el, ok := db.data[op.Key]; ok {
				db.removeElement(el)
				db.stats.Deletes++
			}
		}
	}
}

// Keys returns snapshot of keys (non-blocking for long scans would require different design)
func (db *DB) Keys() []string {
	db.mu.RLock()
	defer db.mu.RUnlock()
	keys := make([]string, 0, len(db.data))
	for k := range db.data {
		keys = append(keys, k)
	}
	return keys
}

// Stats snapshot
func (db *DB) Stats() Stats {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.stats
}

// Close stops the janitor
func (db *DB) Close() {
	close(db.closed)
}

// internals

func (db *DB) removeElement(el *list.Element) {
	kp := el.Value.(*kvPair)
	delete(db.data, kp.key)
	db.stats.Bytes -= uint64(kp.item.size)
	db.lru.Remove(el)
}

func (db *DB) evictLRU() {
	// evict 1 item (could be tuned)
	el := db.lru.Back()
	if el == nil {
		return
	}
	kp := el.Value.(*kvPair)
	delete(db.data, kp.key)
	db.lru.Remove(el)
	db.stats.Evictions++
	db.stats.Bytes -= uint64(kp.item.size)
}

func (db *DB) janitor(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			db.cleanupExpired()
		case <-db.closed:
			return
		}
	}
}

func (db *DB) cleanupExpired() {
	now := time.Now()
	db.mu.Lock()
	defer db.mu.Unlock()

	for el := db.lru.Back(); el != nil; {
		prev := el.Prev()
		kp := el.Value.(*kvPair)
		if kp.item.expiresAt.IsZero() == false && now.After(kp.item.expiresAt) {
			db.removeElement(el)
		}
		el = prev
	}
}
