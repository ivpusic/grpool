# grpool
Simple Goroutine pool

Clients can submit jobs. Dispatcher takes job, and sends it to first available worker.
When worker is done with processing job, will be returned back to worker pool.

Number of workers and Job queue size is configurable.

## Example
```Go
package main

import (
	"fmt"

	"github.com/ivpusic/grpool"
)

func main() {
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
	pool.JobQueue <- job

	// wait until all jobs are done
	pool.Wait()
}

```

## License
*MIT*
