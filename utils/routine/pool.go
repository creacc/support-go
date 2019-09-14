package droutine

import "sync"

type Pool struct {
	sync.Mutex
	workers chan bool
	count   int
}

var DefaultPool = NewPool(2000)

func NewPool(capacity int) *Pool {
	p := &Pool{
		workers: make(chan bool, capacity),
	}
	return p
}

func (p *Pool) Execute(task func()) {
	p.workers <- true
	go func() {
		task()
		<-p.workers
	}()
}
