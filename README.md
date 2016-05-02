# pipeline
Go Pipeline for chained operations

[![Circle CI](https://circleci.com/gh/skidder/pipeline.svg?style=svg)](https://circleci.com/gh/skidder/pipeline) [![GoDoc](https://godoc.org/github.com/skidder/pipeline?status.svg)](https://godoc.org/github.com/skidder/pipeline) [![Go Report Card](https://goreportcard.com/badge/github.com/skidder/pipeline)](https://goreportcard.com/report/github.com/skidder/pipeline)

`go get github.com/skidder/pipeline`

## Overview
Have you ever needed to perform a sequence of distinct operations, possiby in different combinations? Those operations, or pipe segments, could be arranged in a pipeline, with the pipe segments being reusable and possibly ordered differently to compose different pipelines.  This library provides a structure for creating `Pipe` segments and forming them into a reusable `Pipeline`.

## Usage
### Pipes
Your pipe must implement the `Pipe` interface with its `Process` function whose input is a `Data` channel and returns a `Data` channel representing the `Pipe` segment output:

```
-> channel Data -> [Process] -> channel Data ->
```

### Pipelines
A `Pipeline` is constructed from one or more `Pipe` values:

```
-> channel Data -> |                                                                             |-> channel Data ->
                   |                                                                             |
                   -> [Process A] -> channel Data -> [Process B] -> channel Data -> [Process C] ->
```

## Examples
See the [Pass-through Pipe](https://godoc.org/github.com/skidder/pipeline#example-NewPipeline) and [Timing Pipe](https://godoc.org/github.com/skidder/pipeline#example-NewTimingPipe) examples in the Go-Docs.
