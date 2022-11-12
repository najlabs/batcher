package main

import (
	"github.com/najlabs/batcher"
	"log"
)

var beforeEachStepListener = func(ctx *batcher.JobExecutionContext) {
	log.Printf("before step %v\n", ctx.CurrentStep)
}
var afterEachStepListener = func(ctx *batcher.JobExecutionContext) {
	log.Printf("after step %v\n", ctx.CurrentStep)
}

/*
var consoleWriter batcher.ItemWriter[int] = func(i []int) error { fmt.Printf("write: %v\n", i); return nil }
var step1 = batcher.NewStep(dummyStateReader(3), consoleWriter, batcher.WithChunkSize(5), batcher.WithName("Step1"))
var step2 = batcher.NewStep(dummyStateReader(11), consoleWriter, batcher.WithChunkSize(5), batcher.WithName("Step2"))

	func dummyStateReader(index int) batcher.ItemReader[int] {
		return func(buffer []int) (int, error) {
			itemsRead := 0
			for i := 0; i < len(buffer) && index > 0; i++ {
				buffer[i] = index
				index--
				time.Sleep(time.Second * 1)
				itemsRead++
			}
			if index <= 0 {
				return itemsRead, batcher.ErrEnd
			}
			return itemsRead, nil
		}
	}
*/
func exampleWithListener() {
	job := batcher.NewJob("JobX") //, batcher.BeforeEachStep(beforeEachStepListener), batcher.AfterEachStep(afterEachStepListener))
	job.Steps(step1, step2)
	job.Run()
}
func example2() {
	job := batcher.NewJobInline("Job1", step1, step2)
	job.Configure(batcher.BeforeEachStep(beforeEachStepListener), batcher.AfterEachStep(afterEachStepListener))
	job.Run()
}
func main() {
	exampleWithListener()
	example2()
}
