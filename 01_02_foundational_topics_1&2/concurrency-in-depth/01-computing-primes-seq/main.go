package main

import (
	"fmt"
	"math"
	"time"
)

var PRIMES_TILL = 20
var PRIMES_COUNT = 0

func checkPrimes(input int) {
	if input&1 == 0 {
		return
	}
	for i := 3; i <= int(math.Sqrt(float64(input))); i++ {
		if input%i == 0 {
			return
		}
	}
	PRIMES_COUNT++
}

func main() {
	startTime := time.Now()
	for i := 3; i <= PRIMES_TILL; i++ {
		checkPrimes(i)
	}
	fmt.Printf("Time taken = %d \n", time.Since(startTime))
	fmt.Printf("Till %d, Found  %d Primes\n", PRIMES_TILL, PRIMES_COUNT+1)
}
