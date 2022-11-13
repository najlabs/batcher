package main

import (
	"fmt"
	"github.com/najlabs/batcher"
	"log"
	"time"
)

type dummyAction struct {
	index int
}

func (d *dummyAction) Read(buffer []int) (int, error) {
	itemsRead := 0
	for i := 0; i < len(buffer) && d.index > 0; i++ {
		buffer[i] = d.index
		d.index--
		time.Sleep(time.Second * 1)
		itemsRead++
	}
	if d.index <= 0 {
		return itemsRead, batcher.ErrEnd
	}
	return itemsRead, nil
}
func (d *dummyAction) Write(data []int) error {
	fmt.Printf("write: %v\n", data)
	return nil
}

var step1 = batcher.NewStepAction(newDummyAction(3), batcher.WithChunkSize(5), batcher.WithName("Step1"))
var step2 = batcher.NewStepAction(newDummyAction(11), batcher.WithChunkSize(5), batcher.WithName("Step2"))

func newDummyAction(index int) batcher.Step[int] {
	return &dummyAction{index}
}

var beforeEachStepListener = func(ctx *batcher.JobExecutionContext) {
	log.Printf("before step %v\n", ctx.CurrentStep)
}
var afterEachStepListener = func(ctx *batcher.JobExecutionContext) {
	log.Printf("after step %v\n", ctx.CurrentStep)
}
