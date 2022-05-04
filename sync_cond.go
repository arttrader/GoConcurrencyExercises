package main

import (
	"fmt"
	"net"
	"sync"
)

const limit = 5

type Item = int

type Queue struct {
	mu        sync.Mutex
	items     []Item
	itemAdded sync.Cond
}

func NewQueue() *Queue {
	q := new(Queue)
	q.itemAdded.L = &q.mu
	return q
}

func (q *Queue) Put(item Item) {
	q.itemAdded.L.Lock()
	defer q.itemAdded.L.Unlock()
	q.items = append(q.items, item)
	q.itemAdded.Signal() // Signal wakes one goroutine waiting on c, if there is any
}

func (q *Queue) GetMany(n int) []Item {
	q.itemAdded.L.Lock()
	defer q.itemAdded.L.Unlock()
	for len(q.items) < n {
		q.itemAdded.Wait()
	}
	items := q.items[:n:n]
	q.items = q.items[n:]
	return items
}

type Pool struct {
	mu              sync.Mutex
	cond            sync.Cond
	numConns, limit int
	idle            []net.Conn
}

func NewPool(limit int) *Pool {
	p := &Pool{limit: limit}
	p.cond.L = &p.mu
	return p
}

func (p *Pool) Release(c net.Conn) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.idle = append(p.idle, c)
	p.cond.Signal()
}

func (p *Pool) Hijack(c net.Conn) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.numConns--
	p.cond.Signal()
}

func (p *Pool) Acquire() (net.Conn, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for len(p.idle) == 0 &&
		p.numConns >= p.limit {
		p.cond.Wait()
	}

	if len(p.idle) > 0 {
		c := p.idle[len(p.idle)-1]
		p.idle = p.idle[:len(p.idle)-1]
		return c, nil
	}

	c, err := dial()
	if err == nil {
		p.numConns++
	}
	return c, err
}

func dial() (net.Conn, error) {
	p := NewPool(limit)
	conn, err := net.Dial("tcp", "golang.org:80")
	p.idle = append(p.idle, conn)
	return conn, err
}

func main() {
	q := NewQueue()

	var wg sync.WaitGroup
	for n := 10; n > 0; n-- {
		wg.Add(1)
		go func(n int) {
			items := q.GetMany(n)
			fmt.Printf("%2d: %2d\n", n, items)
			wg.Done()
		}(n)
	}

	for i := 0; i < 100; i++ {
		q.Put(i)
	}

	wg.Wait()
}
