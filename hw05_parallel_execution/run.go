package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func checkErrors(errCount int32, m int) error {
	if m > 0 && errCount > int32(m) {
		return ErrErrorsLimitExceeded
	}
	if m <= 0 && errCount > 0 {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		return nil
	}

	taskCh := make(chan Task)
	stop := make(chan struct{})
	var wg sync.WaitGroup
	var errCount int32
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
				err := task()
				if err != nil {
					if m <= 0 {
						once.Do(func() { close(stop) })
					} else if atomic.AddInt32(&errCount, 1) > int32(m) { //nolint:gosec
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
			break
		case taskCh <- task:
		}
	}
	close(taskCh)

	wg.Wait()

	return checkErrors(atomic.LoadInt32(&errCount), m) //nolint:gosec
}
