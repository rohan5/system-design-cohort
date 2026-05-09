package main

import (
	"connection-pooling/db"
	"database/sql"
	"log"
	"sync"
	"time"
)

// build connection pooling in golang using database/sql package and mysql driver
// connection pool should be blocking.
// we will use golang's channel to implement the connection pool.
type conn struct {
	db *sql.DB
}

type cpool struct {
	mu      *sync.Mutex
	channel chan interface{}
	conn    []*conn
	maxConn int
}

func NewCPool(maxConn int) (*cpool, error) {
	pool := &cpool{
		mu:      &sync.Mutex{},
		conn:    make([]*conn, 0, maxConn),
		maxConn: maxConn,
		channel: make(chan interface{}, maxConn),
	}
	for i := 0; i < maxConn; i++ {
		pool.conn = append(pool.conn, &conn{db.New()})
	}
	return pool, nil
}

// func (p *Pool) Get() (*sql.DB, error) {
// }

// func (p *Pool) Put(db *sql.DB) {
// }

func main() {
	log.Println("HEllo World!!!")

	// benchmark connection pool
	benchmarkPool()

	// benchmark non connection pool
	// benchmarkNonPool()
}

func benchmarkPool() {
	startTime := time.Now()
	// new connection pool
	_, err := NewCPool(10)
	if err != nil {
		log.Fatal(err)
	}
	// defer pool.Close()

	// wg := sync.WaitGroup{}
	// wg.Add(100)

	// for i := 0; i < 100; i++ {
	// 	go func() {
	// 		defer wg.Done()
	// 		// conn := Get()

	// 		// _, err := db.Exec("Select SLEEP(0.1);")
	// 		// if err != nil {
	// 		// 	// log.Fatal(err)
	// 		// 	panic(err)
	// 		// }

	// 		// Put(conn)
	// 	}()
	// }
	// wg.Wait()

	log.Printf("** Pool time took = %v **", time.Since(startTime))
}

func benchmarkNonPool() {

	startTime := time.Now()
	count := 200
	wg := sync.WaitGroup{}
	wg.Add(count)
	for i := 0; i < count; i++ {
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
