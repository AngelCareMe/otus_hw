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
	var maxActive int32
	var once sync.Once

	stop := make(chan struct{})

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

				_ = atomic.AddInt32(&maxActive, 1)
				err := task()
				atomic.AddInt32(&maxActive, -1)

				if err != nil {
					if m > 0 {
						if atomic.AddInt32(&errCount, 1) > int32(m) {
							once.Do(func() { close(stop) })
						}
					} else if m <= 0 {
						once.Do(func() { close(stop) })
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
			wg.Wait()
			return ErrErrorsLimitExceeded
		case taskCh <- task:
		}
	}
	close(taskCh)

	wg.Wait()

	if (m > 0 && atomic.LoadInt32(&errCount) > int32(m)) || (m <= 0 && errCount > 0) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
