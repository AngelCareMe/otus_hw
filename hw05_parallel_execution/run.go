package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		return nil
	}

	taskCh := make(chan Task)
	var wg sync.WaitGroup
	var errCount int32

	stop := make(chan struct{})
	var once sync.Once

	worker := func() {
		defer wg.Done()
		for {
			select {
			case <-stop:
				return
			case task, ok := <-taskCh:
				if !ok {
					return
				}
				if err := task(); err != nil {
					if m > 0 && atomic.AddInt32(&errCount, 1) > int32(m) {
						once.Do(func() {
							close(stop)
						})
					} else if m <= 0 {
						once.Do(func() {
							close(stop)
						})
					}
				}
			}
		}
	}

	wg.Add(n)
	for i := 0; i < n; i++ {
		go worker()
	}

	for _, task := range tasks {
		select {
		case <-stop:
			return ErrErrorsLimitExceeded
		default:
			taskCh <- task
		}
	}
	close(taskCh)

	wg.Wait()

	if m > 0 && atomic.LoadInt32(&errCount) > int32(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
