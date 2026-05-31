package main

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

var PRIMES_TILL = 10000000
var PRIMES_COUNT int32 = 0
var BATCHES = 10

func checkPrimes(input int) {
	if input&1 == 0 {
		return
	}
	for i := 3; i <= int(math.Sqrt(float64(input))); i++ {
		if input%i == 0 {
			return
		}
	}
	// PRIMES_COUNT++
	atomic.AddInt32(&PRIMES_COUNT, 1)
}

func main() {
	var wg sync.WaitGroup
	startTime := time.Now()
	batchSize := PRIMES_TILL / BATCHES
	nStart := 3
	for i := 0; i < BATCHES-1; i++ {
		wg.Add(1)
		go processBatch(i, &wg, nStart, nStart+batchSize)
		nStart = nStart + batchSize + 1
	}
	wg.Add(1)
	go processBatch(BATCHES-1, &wg, nStart, PRIMES_TILL)

	wg.Wait()
	fmt.Printf("Time taken = %f seconds\n", time.Since(startTime).Seconds())
	fmt.Printf("Till %d, Found  %d Primes\n", PRIMES_TILL, PRIMES_COUNT+1)
}

func processBatch(batch int, wg *sync.WaitGroup, nStart int, nEnd int) {
	defer wg.Done()
	startTime := time.Now()
	for i := nStart; i <= nEnd; i++ {
		checkPrimes(i)
	}
	fmt.Printf("Batch -> %d [%d, %d] Time : %f \n", batch, nStart, nEnd, time.Since(startTime).Seconds())
}
