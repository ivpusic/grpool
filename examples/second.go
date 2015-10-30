package grpool

import (
	"fmt"
	"runtime"

	"github.com/ivpusic/grpool"
)

func second() {
	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs)

	// number of workers, and size of job queue
	pool := grpool.NewPool(100, 50)
	defer pool.Release()

	// how many jobs we should wait
	pool.WaitCount(10)

	// submit one or more jobs to pool
	for i := 0; i < 10; i++ {
		count := i

		pool.JobQueue <- func() {
			// say that job is done, so we can know how many jobs are finished
			defer pool.JobDone()

			fmt.Printf("hello %d\n", count)
		}
	}

	// wait until we call JobDone for all jobs
	pool.WaitAll()
}
