package in_memory_db

import (
	"fmt"
	"time"
)

// Simple demonstration
func In_mem_db() {
	fmt.Println("Starting demo of in-memory DB")
	db := NewDB(1000, 1*time.Second)
	defer db.Close()

	// Set and Get
	db.Set("name", []byte("anurag"), 0)
	v, err := db.Get("name")
	if err != nil {
		panic(err)
	}
	fmt.Printf("name=%s\n", string(v))

	// TTL
	db.Set("temp", []byte("expiring"), 1)
	v, err = db.Get("temp")
	if err == nil {
		fmt.Printf("temp before expiry=%s\n", string(v))
	}
	fmt.Println("sleeping 2s to let temp expire...")
	time.Sleep(2 * time.Second)
	_, err = db.Get("temp")
	if err != nil {
		fmt.Println("temp expired (expected)")
	}

	// CAS
	db.Set("counter", []byte("1"), 0)
	el := db.data["counter"]
	kp := el.Value.(*kvPair)
	ver := kp.item.ver
	if err := db.CAS("counter", ver, []byte("2"), 0); err != nil {
		fmt.Println("cas failed")
	} else {
		fmt.Println("cas success")
	}

	// Transaction
	tx := NewTx()
	tx.Set("a", []byte("1"), 0)
	tx.Set("b", []byte("2"), 0)
	tx.Delete("name")
	db.Commit(tx)
	if _, err := db.Get("name"); err != nil {
		fmt.Println("name deleted by tx")
	}

	// Stats
	st := db.Stats()
	fmt.Printf("stats: %+v\n", st)
}
