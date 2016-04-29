package pipeline

import (
	"time"
)

// TimingPipe invokes a custom callback function with the amount of time required to run a specific Pipe
type TimingPipe struct {
	timedPipe Pipe
	callback  func(begin time.Time, duration time.Duration)
}

// NewTimingPipe creates a new timing pipe
func NewTimingPipe(timedPipe Pipe, callback func(begin time.Time, duration time.Duration)) Pipe {
	return &TimingPipe{
		timedPipe: timedPipe,
		callback:  callback,
	}
}

func (t *TimingPipe) Process(in chan Data) chan Data {
	out := make(chan Data)
	go func() {
		defer close(out)
		for request := range in {
			begin := time.Now()
			innerPipeline := NewPipeline(t.timedPipe)
			go func() {
				innerPipeline.Enqueue(request)
				innerPipeline.Close()
			}()
			innerPipeline.Dequeue(func(response Data) {
				out <- response
			})
			t.callback(begin, time.Since(begin))
		}
	}()
	return out
}
