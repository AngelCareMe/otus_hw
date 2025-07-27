package hw05parallelexecution

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require" //nolint:depguard
)

func TestConcurrency(t *testing.T) {
	const taskCount = 1000
	const parallelism = 10

	var activeGoroutines int32
	var maxActive int32

	tasks := make([]Task, taskCount)
	for i := range tasks {
		tasks[i] = func() error {
			curr := atomic.AddInt32(&activeGoroutines, 1)
			for {
				currentMax := atomic.LoadInt32(&maxActive)
				if curr <= currentMax {
					break
				}
				if atomic.CompareAndSwapInt32(&maxActive, currentMax, curr) {
					break
				}
			}
			time.Sleep(10 * time.Millisecond)
			atomic.AddInt32(&activeGoroutines, -1)
			return nil
		}
	}

	err := Run(tasks, parallelism, 0)
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		return atomic.LoadInt32(&maxActive) >= int32(parallelism)
	}, time.Second, 10*time.Millisecond, "ожидается, что %d задач выполнялись одновременно", parallelism)
}
