package pipeline

import (
	"time"
)

type timingPipe struct {
	timedPipe Pipe
	callback  func(begin time.Time, duration time.Duration)
}

// NewTimingPipe creates a new timing pipe
func NewTimingPipe(timedPipe Pipe, callback func(begin time.Time, duration time.Duration)) Pipe {
	return &timingPipe{
		timedPipe: timedPipe,
		callback:  callback,
	}
}

func (r *timingPipe) Process(in chan Data) chan Data {
	out := make(chan Data)
	go func() {
		defer close(out)
		for request := range in {
			begin := time.Now()
			innerPipeline := NewPipeline(r.timedPipe)
			go func() {
				innerPipeline.Enqueue(request)
				innerPipeline.Close()
			}()
			innerPipeline.Dequeue(func(response Data) {
				out <- response
			})
			r.callback(begin, time.Since(begin))
		}
	}()
	return out
}
