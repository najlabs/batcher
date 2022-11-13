# batcher
*batcher* is a comprehensive batch processing implementation.

batch processing can be described in its most simple form as reading in large amounts of data, performing some type of calculation or transformation, and writing the result out.
Batch processing are commonly used for:
- automated data migration between application
- convert a list of files to another data format CSV, JSON, Queues, ...
- bulk operation like images resizing, conversion, watermarking
- big data OLTP and ETL processing

## Usage
```go
package main

import (
	"fmt"
	"github.com/najlabs/batcher"
	"time"
)

func main() {
	job := batcher.NewJobInline("Job1", step1, step2)
	job.Run()
}

// create steps
var step1 = batcher.NewStepAction(newDummyAction(3), batcher.WithChunkSize(5), batcher.WithName("Step1"))
var step2 = batcher.NewStepAction(newDummyAction(11), batcher.WithChunkSize(5), batcher.WithName("Step2"))

func newDummyAction(index int) batcher.Step[int] {
	return &dummyAction{index}
}
type dummyAction struct {
	index int
}

func (d *dummyAction) Read(buffer []int) (int, error) {
	fmt.Println("put data in buffer")
	return 0, batcher.ErrEnd
}
func (d *dummyAction) Write(data []int) error {
	fmt.Printf("write: %v\n", data)
	return nil
}
```
for more example see the `./example` folder

## Domain language
### *Job* 
a Job is an entity that encapsulates an entire batch process.
### *Step*
a Step is a domain object that encapsulates an independent, sequential phase of a batch job. It contains all of the information necessary to define and control the actual batch processing.
Every _**Job**_ is composed entirely of _**one or more steps**_. 
and each step can have at least a **Reader**, **Writer** potentially a **mapper**
#### ItemReader
means for providing data from many different types of input like CSV, JSON, Databases, ...
```go
// read data from custom source and put in the buffer
// return the number of data readed and error if any. 
// returning error ErrEnd will mean you ar done with reading
type ItemReader[INPUT any] interface {
    Read([]INPUT) (int, error)
}
```
#### ItemMapper
for Processing or transforming readed data and passing it to Writer
#### ItemWriter
attempts to write out the list of items passed to given destination like CSV files, databases. Queues, ...
```go
type ItemWriter[OUTPUT any] interface {
    Write([]OUTPUT) error
}
```

## Installation

~~~~
go get github.com/najlabs/batcher
~~~~
Or you can manually git clone the repository to
`$(go env GOPATH)/src/github.com/najlabs/batcher`
## create a new Batch
```go
job := batcher.NewJob("Job 1")
job.Steps(step1, step2)
job.Run()
```

## create a new step
```go
// with reader and writer 
var step1 = batcher.NewStep(dummyStateReader(3), consoleWriter)
// with options chunk and name
var step1 = batcher.NewStep(dummyStateReader(3), consoleWriter, batcher.WithChunkSize(5), batcher.WithName("Step1"))
```

