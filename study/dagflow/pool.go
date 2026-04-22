package dagflow

import "sync"

// Pool is a simple goroutine pool that limits concurrency.
type Pool struct {
	taskCh chan func()
	wg     sync.WaitGroup
}

// NewPool creates a pool with the specified number of workers.
func NewPool(workers int) *Pool {
	p := &Pool{
		taskCh: make(chan func(), workers*2),
	}
	for i := 0; i < workers; i++ {
		go p.worker()
	}
	return p
}

func (p *Pool) worker() {
	for task := range p.taskCh {
		task()
	}
}

// Submit adds a task to the pool. Blocks if all workers are busy and the queue is full.
func (p *Pool) Submit(task func()) {
	p.wg.Add(1)
	p.taskCh <- func() {
		defer p.wg.Done()
		task()
	}
}

// Wait blocks until all submitted tasks have completed.
func (p *Pool) Wait() {
	p.wg.Wait()
}

// Shutdown closes the task channel and stops all workers after current tasks finish.
func (p *Pool) Shutdown() {
	close(p.taskCh)
}
