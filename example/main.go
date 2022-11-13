package main

import (
	"github.com/najlabs/batcher"
)

func exampleWithListener() {
	job := batcher.NewJob("JobX", batcher.BeforeEachStep(beforeEachStepListener), batcher.AfterEachStep(afterEachStepListener))
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
