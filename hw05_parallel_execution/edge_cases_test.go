package hw05parallelexecution

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require" //nolint:depguard
)

// Проверка при n <= 0 функция не запускает задачи
func TestRunWithZeroOrNegativeWorkers(t *testing.T) {
	var runTasksCount int32

	tasks := []Task{
		func() error {
			atomic.AddInt32(&runTasksCount, 1)
			return nil
		},
		func() error {
			atomic.AddInt32(&runTasksCount, 1)
			return nil
		},
	}

	err := Run(tasks, 0, 10)

	require.NoError(t, err)
	require.Zero(t, atomic.LoadInt32(&runTasksCount), "tasks should not be executed when n <= 0")
}

// Проверка не запускается больше n воркеров
func TestRunDoesNotExceedWorkerLimit(t *testing.T) {
	const expectedWorkers = 5
	var startedWorkers int32
	var maxWorkers int32

	tasks := make([]Task, 100)
	for i := range tasks {
		tasks[i] = func() error {
			curr := atomic.AddInt32(&startedWorkers, 1)
			for {
				currentMax := atomic.LoadInt32(&maxWorkers)
				if curr <= currentMax {
					break
				}
				if atomic.CompareAndSwapInt32(&maxWorkers, currentMax, curr) {
					break
				}
			}
			time.Sleep(1 * time.Millisecond)
			atomic.AddInt32(&startedWorkers, -1)
			return nil
		}
	}

	err := Run(tasks, expectedWorkers, 10)
	require.NoError(t, err)
	require.LessOrEqual(t, atomic.LoadInt32(&maxWorkers), int32(expectedWorkers), "should not exceed worker limit")
}
