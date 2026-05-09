package main

import (
	"connection-pooling/db"
	"database/sql"
	"fmt"
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

func (p *cpool) Get() (*conn, error) {
	conn := p.conn[0]
	p.conn = p.conn[1:]
	if conn == nil {
		return nil, fmt.Errorf("no connection available")
	}
	return conn, nil
}

func (p *cpool) Put(conn *conn) {
	p.conn = append(p.conn, conn)
}

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
	cpool, err := NewCPool(10)
	if err != nil {
		log.Fatal(err)
	}

	// conn, err := cpool.Get()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// res, err := conn.db.Exec("Select SLEEP(0.1);")
	// if err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	log.Println("Result: ", res)
	// }
	// cpool.Put(conn)
	goRoutinesCount := 11
	wg := sync.WaitGroup{}
	wg.Add(goRoutinesCount)

	for i := 0; i < goRoutinesCount; i++ {
		go func() {
			defer wg.Done()
			conn, err := cpool.Get()
			if err != nil {
				log.Fatal(err)
			}
			res, err := conn.db.Exec("Select SLEEP(0.1);")
			if err != nil {
				log.Fatal(err)
			} else {
				log.Println("Result: ", res)
			}
			cpool.Put(conn)
		}()
	}
	wg.Wait()

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
