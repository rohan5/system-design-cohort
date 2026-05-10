package main

import (
	"connection-pooling/db"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"
)

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
		pool.channel <- nil
	}
	return pool, nil
}

func (p *cpool) Get() (*conn, error) {
	<-p.channel
	p.mu.Lock()
	conn := p.conn[0]
	p.conn = p.conn[1:]
	if conn == nil {
		return nil, fmt.Errorf("no connection available")
	}
	p.mu.Unlock()
	return conn, nil
}

func (p *cpool) Put(conn *conn) {
	p.mu.Lock()
	p.conn = append(p.conn, conn)
	p.mu.Unlock()

	p.channel <- nil
}

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
	cpool, err := NewCPool(10)
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
			_, err = conn.db.Exec("Select SLEEP(0.1);")
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
