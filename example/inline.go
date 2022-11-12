package main

import (
	"fmt"
	"github.com/najlabs/batcher"
	"time"
)

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

var consoleWriter batcher.ItemWriter[int] = func(i []int) error { fmt.Printf("write: %v\n", i); return nil }

/*func main() {
	job := batcher.NewJobInline("Job1", step1, step2)
	job.Run()
}*/
