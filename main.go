package main

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/najlabs/batcher/core"
)

func main() {
	// csvInputReader := csvreader.Read("").Mapper()
	// consoleWriter := consolewriter.Read("").Mapper()
	// step1 := core.NewStepBuilder[int, string]().Reader(csvInputReader).Writer(consoleWriter)
	// var simpleReader core.ItemReader[int] = func() (int, error) { return 1, nil }
	// var simpleReader1 core.ItemReader[int] = simpleStateReader()
	// var simpleProcessor core.ItemProcessor[int, int] = func(i int) (int, error) { return i + 1, nil }
	// var simpleProcessor2 core.ItemProcessor[int, string] = func(i int) (string, error) { return fmt.Sprintf("%d", i+1), nil }
	// var consoleWriter core.ItemWriter[int] = func(i []int) error { fmt.Printf("write: %v\n", i); return nil }
	// var consoleWriter1 core.ItemWriter[string] = func(i []string) error { fmt.Printf("write: %v\n", i); return nil }
	// var simpleReader2 core.ItemReader[int] = func(ctx *core.BatchContext) (int, error) { return 0, nil }
	step1 := core.NewSimpleStep(read, write).WithChunkSizeOf(2)
	step2 := core.NewSimpleStep(readString, writeString).WithChunkSizeOf(2)
	// step2 := core.NewStep(simpleStateReader(), simpleProcessor, consoleWriter)
	// step3 := core.NewStep(simpleStateReader(), simpleProcessor2, consoleWriter1)
	core.NewJob(step1, step2).Run()
}

/*
	 func simpleStateReader() core.ItemReader[int] {
		index := 5
		return func() (int, error) {
			time.Sleep(time.Second * 1)
			index--
			if index < 0 {
				return 0, nil
			}
			return index, nil
		}
	}
*/
func readString(ctx *core.ReaderContext[string]) error {
	id := uuid.NewString()
	time.Sleep(time.Second * 1)
	return ctx.AddNextItem(id)
}
func writeString(ctx *core.WriterContext[string]) error {
	for _, value := range ctx.ItemsToWrite() {
		log.Println(value)
	}
	return ctx.Done()
}
func read(ctx *core.ReaderContext[int]) error {
	time.Sleep(time.Second * 1)
	index := ctx.GetOrCreate("index", 9).(int)
	if index > 0 {
		index--
		ctx.UpdateIfExists("index", index)
		return ctx.AddNextItem(index)
	}
	return ctx.Done()
}

func write(ctx *core.WriterContext[int]) error {
	log.Print("Writer ....", ctx.ItemsToWrite())
	for _, value := range ctx.ItemsToWrite() {
		time.Sleep(time.Second * 1)
		log.Println("--> write ", value)
	}
	return ctx.Done()
}
