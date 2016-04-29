package pipeline

import (
	"os"
	"testing"
	"time"
)

type TimingPipeTestCase struct {
	data     Data
	callback func(begin time.Time, duration time.Duration)
}

func TestTimingPipeProcess(t *testing.T) {
	tests := make(map[string]*TimingPipeTestCase)
	tests["Success"] = &TimingPipeTestCase{
		data:     Data{Payload: "foo"},
		callback: func(begin time.Time, duration time.Duration) {},
	}

	for name, tc := range tests {
		pipes := NewPipeline(NewTimingPipe(&PassThroughPipe{}, tc.callback))
		go func() {
			pipes.Enqueue(tc.data)
			pipes.Close()
		}()

		var pipeOutput Data
		pipes.Dequeue(func(data Data) {
			pipeOutput = data
		})
		if pipeOutput.Payload.(string) != "foo" {
			t.Error("%s: pipeline paylaod string had incorrect value: %s", name, pipeOutput.Payload.(string))
		}
	}

}

func ExampleTimingPipe() {
	os.Stdout.Sync()
	timingCallbackFoo := func(begin time.Time, duration time.Duration) {
		os.Stdout.WriteString("foo")
	}
	timingCallbackBar := func(begin time.Time, duration time.Duration) {
		os.Stdout.WriteString("bar")
	}

	pipes := NewPipeline(
		NewTimingPipe(&PassThroughPipe{}, timingCallbackFoo),
		NewTimingPipe(&PassThroughPipe{}, timingCallbackBar),
	)
	go func() {
		pipes.Enqueue(Data{})
		pipes.Close()
	}()

	var pipeOutput Data
	pipes.Dequeue(func(data Data) {
		pipeOutput = data
	})
	// Output: foobar
}
