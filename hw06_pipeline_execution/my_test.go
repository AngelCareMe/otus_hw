package hw06pipelineexecution

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmptyPipeline(t *testing.T) {
	in := make(Bi)
	go func() {
		close(in)
	}()
	result := make([]interface{}, 0)
	for v := range ExecutePipeline(in, nil, []Stage{}...) {
		result = append(result, v)
	}
	require.Empty(t, result, "Empty pipeline should produce no output")
}

func TestSingleStage(t *testing.T) {
	in := make(Bi)
	data := []int{1, 2, 3}
	go func() {
		for _, v := range data {
			in <- v
		}
		close(in)
	}()
	stage := func(in In) Out {
		out := make(Bi)
		go func() {
			defer close(out)
			for v := range in {
				out <- v.(int) + 1
			}
		}()
		return out
	}
	result := make([]interface{}, 0)
	for v := range ExecutePipeline(in, nil, stage) {
		result = append(result, v)
	}
	require.Equal(t, []interface{}{2, 3, 4}, result)
}
