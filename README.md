# grpool
[![Build Status](https://travis-ci.org/ivpusic/grpool.svg?branch=master)](https://travis-ci.org/ivpusic/grpool)

Lightweight Goroutine pool

Clients can submit jobs. Dispatcher takes job, and sends it to first available worker.
When worker is done with processing job, will be returned back to worker pool.

Number of workers and Job queue size is configurable.

## Docs
https://godoc.org/github.com/ivpusic/grpool

## Installation
```
go get github.com/ivpusic/grpool
```

## Simple example
```Go
package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/ivpusic/grpool"
)

func main() {
	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs)

	// number of workers, and size of job queue
	pool := grpool.NewPool(10, 50)

	job := grpool.Job{
		// define function which will be called on worker
		Fn: func(arg interface{}) {
			msg := arg.(string)
			fmt.Printf("hello %s\n", msg)
		},
		// provide arguments
		Arg: "world",
	}

	// submit one or more jobs to pool
	for i := 0; i < 10; i++ {
		pool.JobQueue <- job
	}

	// dummy wait until jobs are finished
	time.Sleep(1 * time.Second)
}
```

## Example with waiting jobs to finish
```Go
package main

import (
	"fmt"
	"runtime"

	"github.com/ivpusic/grpool"
)

func main() {
	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs)

	// number of workers, and size of job queue
	pool := grpool.NewPool(10, 50)

	// how many jobs we should wait
	pool.WaitCount(10)

	job := grpool.Job{
		// define function which will be called on worker
		Fn: func(arg interface{}) {
			// say that job is done, so we can know how many jobs are finished
			defer pool.JobDone()

			msg := arg.(string)
			fmt.Printf("hello %s\n", msg)
		},
		// provide arguments
		Arg: "world",
	}

	// submit one or more jobs to pool
	for i := 0; i < 10; i++ {
		pool.JobQueue <- job
	}

	// wait until we call JobDone for all jobs
	pool.WaitAll()
}
```

## License
*MIT*
