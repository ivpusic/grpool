package grpool

import (
	"fmt"
	"runtime"
	"time"

	"github.com/ivpusic/grpool"
)

func first() {
	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs)

	// number of workers, and size of job queue
	pool := grpool.NewPool(100, 50)

	// release resources used by pool
	defer pool.Release()

	// submit one or more jobs to pool
	for i := 0; i < 10; i++ {
		count := i

		pool.JobQueue <- func() {
			fmt.Printf("I am worker! Number %d\n", count)
		}
	}

	// dummy wait until jobs are finished
	time.Sleep(1 * time.Second)
}
