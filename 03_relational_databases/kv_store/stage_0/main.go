package main

import (
	"kv_store/stage_0/store"
	"sync"
)

func main() {
	// This main function is intentionally left empty.
	// The actual logic for creating the KV store is in the create_kv_store.go script.

	// store.Put(`HI`, `there?`, nil)
	var wg sync.WaitGroup
	wg.Add(50)
	for i := 0; i < 50; i++ {
		go func(value int) {
			store.Put(`HI`, value, nil)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
