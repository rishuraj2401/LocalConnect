package worker

import (
	"context"
	"sync"
)

type Job func(ctx context.Context)

type Pool struct {
	jobs   chan Job
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func NewPool(buffer int, workers int) *Pool {
	ctx, cancel := context.WithCancel(context.Background())
	pool := &Pool{
		jobs:   make(chan Job, buffer),
		ctx:    ctx,
		cancel: cancel,
	}
	for i := 0; i < workers; i++ {
		pool.wg.Add(1)
		go pool.worker()
	}
	return pool
}

func (p *Pool) worker() {
	defer p.wg.Done()
	for job := range p.jobs {
		job(p.ctx)
	}
}

func (p *Pool) Submit(job Job) {
	select {
	case p.jobs <- job:
	default:
		go job(p.ctx)
	}
}

func (p *Pool) Shutdown() {
	p.cancel()
	close(p.jobs)
	p.wg.Wait()
}
