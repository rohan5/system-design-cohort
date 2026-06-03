package main

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

var PRIMES_TILL = 100000000
var PRIMES_COUNT int32 = 0
var CONCURRENCY = 10
var CURR_NUMBER int32 = 2

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

	for i := 0; i < CONCURRENCY; i++ {
		wg.Add(1)
		go doWork(i, &wg)
	}

	wg.Wait()
	fmt.Printf("Time taken = %f seconds\n", time.Since(startTime).Seconds())
	fmt.Printf("Till %d, Found  %d Primes\n", PRIMES_TILL, PRIMES_COUNT+1)
}

func doWork(thread int, wg *sync.WaitGroup) {
	defer wg.Done()
	startTime := time.Now()
	for {
		x := atomic.AddInt32(&CURR_NUMBER, 1)
		if x > int32(PRIMES_TILL) {
			break
		}
		checkPrimes(int(x))
	}

	fmt.Printf("Batch -> %d Time : %f \n", thread, time.Since(startTime).Seconds())

}
