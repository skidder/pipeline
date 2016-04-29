package pipeline

import (
	"io/ioutil"
	"os"
	"strings"
)

// Data to be passed between pipe segments composing a pipeline
type Data struct {
	Payload  interface{}
	tempDirs []string
}

// Pipe is a segment in a pipeline that can process a given map of job attributes
type Pipe interface {
	Process(in chan Data) chan Data
}

// Pipeline composed of channels for head & tail
type Pipeline struct {
	head chan Data
	tail chan Data
}

// NewPipeline returns a new pipeline composed of the set of supplied pipes
func NewPipeline(pipes ...Pipe) Pipeline {
	head := make(chan Data)
	var nextChan chan Data
	for _, pipe := range pipes {
		if nextChan == nil {
			nextChan = pipe.Process(head)
		} else {
			nextChan = pipe.Process(nextChan)
		}
	}
	return Pipeline{head: head, tail: nextChan}
}

// Enqueue an item in the pipeline
func (p *Pipeline) Enqueue(item Data) {
	p.head <- item
}

// Dequeue an item from the pipeline
func (p *Pipeline) Dequeue(handler func(Data)) {
	for i := range p.tail {
		handler(i)
	}
}

// Close the pipeline
func (p *Pipeline) Close() {
	close(p.head)
}

// CreateTempDir makes a temporary directory associated tied to this pipeline data
func (d *Data) CreateTempDir(tmpDir string, prefix string) (string, error) {
	if d.tempDirs == nil {
		d.tempDirs = make([]string, 0)
	}

	tempDir, err := ioutil.TempDir(strings.TrimSpace(tmpDir), prefix)
	if err != nil {
		return "", err
	}

	d.tempDirs = append(d.tempDirs, tempDir)
	return tempDir, nil
}

// DeleteTempDirs removes temporary directories created with CreateTempDir
func (d *Data) DeleteTempDirs() {
	if d.tempDirs == nil {
		return
	}

	for _, f := range d.tempDirs {
		os.RemoveAll(f)
	}
	d.tempDirs = nil
	// Output: foo
}
