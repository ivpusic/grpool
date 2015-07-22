package grpool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWorker(t *testing.T) {
	pool := make(chan worker)
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

func TestNewDispatcher(t *testing.T) {
	jobQueue := make(chan Job, 20)
	workerPool := make(chan worker, 10)

	d := newDispatcher(workerPool, jobQueue)
	assert.NotNil(t, d)

	counter := 0
	iterations := 10000

	for i := 0; i < iterations; i++ {
		job := Job{
			Fn: func(arg interface{}) {
				val := arg.(int)
				counter += val
				assert.Equal(t, 1, val)
			},
			Arg: 1,
		}

		d.jobQueue <- job
	}

	d.stop()
	assert.Equal(t, iterations, counter)
}

func TestNewPool(t *testing.T) {
	pool := NewPool(10, 100)

	counter := 0
	iterations := 1000000

	for i := 0; i < iterations; i++ {
		job := Job{
			Fn: func(arg interface{}) {
				val := arg.(int)
				counter += val
				assert.Equal(t, 1, val)
			},
			Arg: 1,
		}

		pool.JobQueue <- job
	}

	pool.Wait()
	assert.Equal(t, iterations, counter)
}
