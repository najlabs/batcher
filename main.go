package main

import (
	"fmt"
	"time"

	"github.com/najlabs/batcher/core"
)

func main() {
	// csvInputReader := csvreader.Read("").Mapper()
	// consoleWriter := consolewriter.Read("").Mapper()
	// step1 := core.NewStepBuilder[int, string]().Reader(csvInputReader).Writer(consoleWriter)
	// var simpleReader core.ItemReader[int] = func() (int, error) { return 1, nil }
	var simpleReader1 core.ItemReader[int] = simpleStateReader()
	var simpleProcessor core.ItemProcessor[int, int] = func(i int) (int, error) { return i + 1, nil }
	var simpleProcessor2 core.ItemProcessor[int, string] = func(i int) (string, error) { return fmt.Sprintf("%d", i+1), nil }
	var consoleWriter core.ItemWriter[int] = func(i int) error { fmt.Printf("write: %v\n", i); return nil }
	var consoleWriter1 core.ItemWriter[string] = func(i string) error { fmt.Printf("write: %v\n", i); return nil }

	step1 := core.NewSimpleStep(simpleStateReader(), consoleWriter)
	step2 := core.NewStep(simpleReader1, simpleProcessor, consoleWriter)
	step3 := core.NewStep(simpleReader1, simpleProcessor2, consoleWriter1)
	core.NewJob(step1, step2, step3).Run()

}

func simpleStateReader() core.ItemReader[int] {
	index := 5
	return func() (int, error) {
		time.Sleep(time.Second * 1)
		if index < 0 {
			return 0, nil //core.ErrStepNoValueToProcess
		}
		index--
		return index, nil
	}
}
