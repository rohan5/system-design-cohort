package dbConnPool

import (
	"connection-pooling/db"
	"fmt"
	"sync"

	"database/sql"
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
	defer p.mu.Unlock()
	conn := p.conn[0]
	p.conn = p.conn[1:]
	if conn == nil {
		return nil, fmt.Errorf("no connection available")
	}
	return conn, nil
}

func (p *cpool) Put(conn *conn) {
	p.mu.Lock()
	p.conn = append(p.conn, conn)
	p.mu.Unlock()

	p.channel <- nil
}

func (c *conn) Exec(query string, args ...any) (sql.Result, error) {
	return c.db.Exec(query, args...)
}
