package grpool

// Gorouting instance which can accept client jobs
type worker struct {
	workerPool chan worker
	jobChannel chan Job
}

func (w worker) start() {
	go func() {
		for {
			// worker free, add it to pool
			w.workerPool <- w

			select {
			case job := <-w.jobChannel:
				if job.stop {
					return
				}

				job.Fn(job.Arg)
			}
		}
	}()
}

func newWorker(pool chan worker) worker {
	return worker{
		workerPool: pool,
		jobChannel: make(chan Job),
	}
}

// Accepts jobs from clients, and waits for first free worker to deliver job
type dispatcher struct {
	workerPool chan worker
	jobQueue   chan Job
}

func (d dispatcher) dispatch() {
	for {
		select {
		case job := <-d.jobQueue:
			worker := <-d.workerPool
			worker.jobChannel <- job
		}
	}
}

func (d dispatcher) stop() {
	for i := 0; i < cap(d.workerPool); i++ {
		worker := <-d.workerPool
		worker.jobChannel <- Job{stop: true}
	}
}

func newDispatcher(workerPool chan worker, jobQueue chan Job) dispatcher {
	d := dispatcher{workerPool: workerPool, jobQueue: jobQueue}

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
type Job struct {
	stop bool
	Fn   func(interface{})
	Arg  interface{}
}

type Pool struct {
	JobQueue   chan Job
	dispatcher dispatcher
}

// Will make pool of gorouting workers.
// numWorkers - how many workers will be created for this pool
// queueLen - how many jobs can we accept until we block
//
// Returned object contains JobQueue reference, which you can use to send job to pool.
func NewPool(numWorkers int, jobQueueLen int) *Pool {
	jobQueue := make(chan Job, jobQueueLen)
	workerPool := make(chan worker, numWorkers)

	return &Pool{
		JobQueue:   jobQueue,
		dispatcher: newDispatcher(workerPool, jobQueue),
	}
}

// Will send special kind of job to all workers to exit.
func (p *Pool) Wait() {
	p.dispatcher.stop()
}
