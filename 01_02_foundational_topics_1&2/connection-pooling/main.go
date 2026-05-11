package main

import (
	"connection-pooling/db"
	"connection-pooling/dbConnPool"
	"log"
	"sync"
	"time"
)

func main() {
	log.Println("HEllo World!!!")

	// benchmark connection pool
	gotoutinesCount := 500
	log.Printf("** Benchmarking with %d goroutines **", gotoutinesCount)

	benchmarkPool(gotoutinesCount)
	// benchmarkNonPool(gotoutinesCount)
}

func benchmarkPool(goRoutinesCount int) {
	// new connection pool
	cpool, err := dbConnPool.NewCPool(10)
	if err != nil {
		log.Fatal(err)
	}

	startTime := time.Now()
	// goRoutinesCount := 11
	wg := sync.WaitGroup{}
	wg.Add(goRoutinesCount)

	for i := 0; i < goRoutinesCount; i++ {
		go func() {
			defer wg.Done()
			conn, err := cpool.Get()
			if err != nil {
				log.Fatal(err)
			}
			_, err = conn.Exec("Select SLEEP(0.1);")
			if err != nil {
				log.Fatal(err)
			}
			cpool.Put(conn)
		}()
	}
	wg.Wait()

	log.Printf("** Pool time took = %v **", time.Since(startTime))
}

func benchmarkNonPool(goRoutinesCount int) {

	startTime := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(goRoutinesCount)
	for i := 0; i < goRoutinesCount; i++ {
		go func() {
			defer wg.Done()
			db := db.New()
			_, err := db.Exec("Select SLEEP(0.1);")
			if err != nil {
				// log.Fatal(err)
				panic(err)
			}
			db.Close()
		}()
	}
	wg.Wait()
	log.Printf("** NonPool time took = %v **", time.Since(startTime))
}
