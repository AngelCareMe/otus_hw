package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	in = wrapWithDone(in, done)
	for _, stage := range stages {
		in = wrapWithDone(stage(in), done)
	}
	return in
}

func wrapWithDone(in In, done In) Out {
	out := make(chan interface{})
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				// дочитываем канал, чтобы upstream не заблокировался
				for range in {
				}
				return
			case val, ok := <-in:
				if !ok {
					return
				}
				select {
				case <-done:
					for range in {
					}
					return
				case out <- val:
				}
			}
		}
	}()
	return out
}
