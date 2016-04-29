package pipeline

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

var (
	tempDirPrefix = "pipeline-test"
)

type CreateTempDirTestCase struct {
	tmpDir        string
	expectError   bool
	expectDirName bool
	setupFunc     func()
	teardownFunc  func()
}

func TestPipeline(t *testing.T) {
	pipes := NewPipeline(&PassThroughPipe{}, &PassThroughPipe{}, &PassThroughPipe{})

	go func() {
		pipes.Enqueue(Data{Payload: "foo"})
		pipes.Close()
	}()

	var pipeOutput Data
	pipes.Dequeue(func(data Data) {
		pipeOutput = data
	})
	if pipeOutput.Payload.(string) != "foo" {
		t.Fail()
	}
}

func TestCreateTempDir(t *testing.T) {
	var tmpBaseDir string
	tests := map[string]*CreateTempDirTestCase{
		"Empty Temp Base Directory": &CreateTempDirTestCase{
			expectError:   false,
			expectDirName: true,
			setupFunc:     func() { tmpBaseDir = "" },
			teardownFunc:  func() {},
		},
		"Whitespace Temp Base Directory": &CreateTempDirTestCase{
			tmpDir:        "  ",
			expectError:   false,
			expectDirName: true,
			setupFunc:     func() { tmpBaseDir = " " },
			teardownFunc:  func() {},
		},
		"Permission Error on Temp Base Directory": &CreateTempDirTestCase{
			expectError:   true,
			expectDirName: false,
			setupFunc: func() {
				var err error
				tmpBaseDir, err = ioutil.TempDir("", "tempdirtest")
				if err != nil {
					panic(err.Error())
				}
				os.Chmod(tmpBaseDir, os.FileMode(0444))
			},
			teardownFunc: func() {
				os.RemoveAll(tmpBaseDir)
			},
		},
		"Success": &CreateTempDirTestCase{
			expectError:   false,
			expectDirName: true,
			setupFunc: func() {
				var err error
				tmpBaseDir, err = ioutil.TempDir("", "tempdirtest")
				if err != nil {
					panic(err.Error())
				}
			},
			teardownFunc: func() {
				os.RemoveAll(tmpBaseDir)
			},
		},
	}

	for name, tc := range tests {
		data := &Data{}
		defer data.DeleteTempDirs()
		defer tc.teardownFunc()

		tc.setupFunc()

		path, err := data.CreateTempDir(tmpBaseDir, tempDirPrefix)

		if tc.expectDirName {
			if "" == path {
				t.Errorf("%s: Expected non-empty directory name but was empty", name)
			}
		} else if "" != path {
			t.Errorf("%s: Expected empty directory name but was non-empty", name)
		}

		if tc.expectError {
			if err == nil {
				t.Errorf("%s: Expected error but there was none", name)
			}
		} else {
			if err != nil {
				t.Errorf("%s: Expected no error but an error was reported", name)
			}
			if 1 != len(data.tempDirs) {
				t.Errorf("%s: Expected the temp-files collection to have one entry", name)
			}
		}
	}
}

func TestDeleteTempDirs(t *testing.T) {
	data := &Data{}
	data.CreateTempDir("", tempDirPrefix)
	data.DeleteTempDirs()

	if data.tempDirs != nil {
		t.Error("Expected the temp-files collection to be nil")
	}
}

func TestDeleteTempDirsEmpty(t *testing.T) {
	data := &Data{}
	data.DeleteTempDirs()

	if data.tempDirs != nil {
		t.Error("Expected the temp-files collection to be nil")
	}
}

func ExamplePipeline() {
	pipes := NewPipeline(&PassThroughPipe{}, &PassThroughPipe{}, &PassThroughPipe{})

	go func() {
		pipes.Enqueue(Data{Payload: "foo"})
		pipes.Close()
	}()

	var pipeOutput Data
	pipes.Dequeue(func(data Data) {
		pipeOutput = data
	})
	fmt.Println(pipeOutput.Payload.(string))
	// Output: foo
}
