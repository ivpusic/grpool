package grpool

import (
	"log"
	"sync"
)

// Gorouting instance which can accept client jobs
type worker struct {
	workerPool chan *worker
	jobChannel chan Job
}

func (w *worker) start() {
	go func() {
		var job Job

		// This defer function will try to catches a crash
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
				if job.RecoverFn != nil {
					job.RecoverFn(r)
				}
			}
		}()

		for {
			// worker free, add it to pool
			w.workerPool <- w

			select {
			case job = <-w.jobChannel:
				job.Fn(job.Arg)
			}
		}
	}()
}

func newWorker(pool chan *worker) *worker {
	return &worker{
		workerPool: pool,
		jobChannel: make(chan Job),
	}
}

// Accepts jobs from clients, and waits for first free worker to deliver job
type dispatcher struct {
	workerPool chan *worker
	jobQueue   chan Job
}

func (d *dispatcher) dispatch() {
	for {
		select {
		case job := <-d.jobQueue:
			worker := <-d.workerPool
			worker.jobChannel <- job
		}
	}
}

func newDispatcher(workerPool chan *worker, jobQueue chan Job) dispatcher {
	d := dispatcher{
		workerPool: workerPool,
		jobQueue:   jobQueue,
	}

	for i := 0; i < cap(d.workerPool); i++ {
		worker := newWorker(d.workerPool)
		worker.start()
	}

	go d.dispatch()
	return d
}

// Represents user request.
// User has to provide function and optional arguments.
// Job will be executed in first free goroutine
// You can supply custom recover function for a panic.
type Job struct {
	Fn        func(interface{})
	Arg       interface{}
	RecoverFn func(interface{})
}

type Pool struct {
	JobQueue   chan Job
	dispatcher dispatcher
	wg         sync.WaitGroup
}

// Will make pool of gorouting workers.
// numWorkers - how many workers will be created for this pool
// queueLen - how many jobs can we accept until we block
//
// Returned object contains JobQueue reference, which you can use to send job to pool.
func NewPool(numWorkers int, jobQueueLen int) *Pool {
	jobQueue := make(chan Job, jobQueueLen)
	workerPool := make(chan *worker, numWorkers)

	pool := &Pool{
		JobQueue:   jobQueue,
		dispatcher: newDispatcher(workerPool, jobQueue),
	}

	return pool
}

// In case you are using WaitAll fn, you should call this method
// every time your job is done.
//
// If you are not using WaitAll then we assume you have your own way of synchronizing.
func (p *Pool) JobDone() {
	p.wg.Done()
}

// How many jobs we should wait when calling WaitAll.
// It is using WaitGroup Add/Done/Wait
func (p *Pool) WaitCount(count int) {
	p.wg.Add(count)
}

// Will wait for all jobs to finish.
func (p *Pool) WaitAll() {
	p.wg.Wait()
}
