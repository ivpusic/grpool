package grpool

import (
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	println("using MAXPROC")
	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs)
}

func TestNewWorker(t *testing.T) {
	pool := make(chan *worker)
	worker := newWorker(pool)
	worker.start()
	assert.NotNil(t, worker)

	worker = <-pool
	assert.NotNil(t, worker, "Worker should register itself to the pool")

	called := false
	done := make(chan bool)

	job := Job{
		Fn: func(_ interface{}) {
			called = true
			done <- true
		},
	}

	worker.jobChannel <- job
	<-done
	assert.Equal(t, true, called)
}

func TestNewPool(t *testing.T) {
	pool := NewPool(1000, 10000)
	defer pool.Release()

	iterations := 1000000
	pool.WaitCount(iterations)
	var counter uint64 = 0

	for i := 0; i < iterations; i++ {
		job := Job{
			Fn: func(arg interface{}) {
				defer pool.JobDone()

				val := arg.(uint64)
				atomic.AddUint64(&counter, val)
				assert.Equal(t, uint64(1), val)
			},
			Arg: uint64(1),
		}

		pool.JobQueue <- job
	}

	pool.WaitAll()

	counterFinal := atomic.LoadUint64(&counter)
	assert.Equal(t, uint64(iterations), counterFinal)
}

func TestRelease(t *testing.T) {
	grNum := runtime.NumGoroutine()
	pool := NewPool(5, 10)
	defer func() {
		pool.Release()

		// give some time for all goroutines to quit
		assert.Equal(t, grNum, runtime.NumGoroutine(), "All goroutines should be released after Release() call")
	}()

	pool.WaitCount(1000)

	for i := 0; i < 1000; i++ {
		job := Job{
			Fn: func(arg interface{}) {
				defer pool.JobDone()
			},
		}

		pool.JobQueue <- job
	}

	pool.WaitAll()
}

func TestCustomRecover(t *testing.T) {
	pool := NewPool(1, 1)
	defer pool.Release()

	pool.WaitCount(1)
	var count int = 1

	job := Job{
		Fn: func(arg interface{}) {
			defer pool.JobDone()
			panic("Capture a custom panic!")
		},
		Arg: uint64(1),
		RecoverFn: func(r interface{}) {
			count++
		},
	}

	pool.JobQueue <- job

	pool.WaitAll()

	time.Sleep(1 * time.Second)
	assert.Equal(t, 2, count)
}

func TestRecover(t *testing.T) {
	pool := NewPool(1, 1)
	defer pool.Release()

	pool.WaitCount(1)

	job := Job{
		Fn: func(arg interface{}) {
			defer pool.JobDone()
			panic("Capture a default panic!")
		},
		Arg: uint64(1),
	}

	pool.JobQueue <- job

	pool.WaitAll()
}
