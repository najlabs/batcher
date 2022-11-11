package core

import (
	"errors"
	"log"

	"github.com/google/uuid"
)

var ErrIncorrectMapper = errors.New("incorrect mapper")

const DefaultStepChunkSize = 10

type Item map[string]any

type ItemReader[I any] func(ctx *ReaderContext[I]) error
type ItemWriter[O any] func(ctx *WriterContext[O]) error
type ItemProcessor[I any, O any] func(ctx *BatchContext, i I) (O, error)

// Step is a building art of a job
// https://docs.spring.io/spring-batch/docs/current/reference/html/domain.html#domainLanguageOfBatch
type Step interface {
	Execute() error
}
type DefaultStep[I any, O any] struct {
	name      string
	reader    ItemReader[I]
	writer    ItemWriter[O]
	processor ItemProcessor[I, O]
	chunkSize int
}

func (step *DefaultStep[I, O]) commitBuffer(readerCtx *ReaderContext[I]) (err error) {
	log.Println("commit reader ")
	println(readerCtx.bufferedItems)
	if len(readerCtx.bufferedItems) == 0 {
		return
	}
	writerCtx := NewWriterContext[O](readerCtx.store)

	if writerCtx.itemToWrite, err = step.mapIO(readerCtx.bufferedItems); err == nil {
		err = step.writer(writerCtx)
		readerCtx.ResetBuffer()
	}
	return
}
func (step *DefaultStep[I, O]) mapIO(inputsData []I) ([]O, error) {
	outputs := make([]O, 0)
	log.Printf("inputsData: %#v, outputs %#v", inputsData, outputs)
	for _, inputItem := range inputsData {
		outputItem, ok := any(inputItem).(O)
		if !ok {
			return outputs, ErrIncorrectMapper
		}
		outputs = append(outputs, outputItem)
	}
	return outputs, nil
}
func (step *DefaultStep[I, O]) Execute() (err error) {
	log.Printf("Executing step %v\n", step.name)
	readerCtx := NewReaderContext[I](step.chunkSize)
looplevel1:
	for {
		switch err = step.reader(readerCtx); err {
		case readerCtx.Cancel():
			log.Println("cancel reader ")
			break looplevel1
		case readerCtx.Continue():
			continue
		case readerCtx.CommitBuffered():
			err = step.commitBuffer(readerCtx)
		case readerCtx.Done():
			log.Println("Done reader ")
			_ = step.commitBuffer(readerCtx)
			err = nil
			break looplevel1
		case readerCtx.Skip():
			log.Println("continue reader ")
			continue looplevel1
		default:
			log.Fatal(err)
		}
	}
	log.Printf("Finishing step %v\n", step.name)
	return
}

/* func (step *DefaultStep[I, O]) Execute() (err error) {
	log.Printf("Executing step %v\n", step.name)
	for {
		chunckedItems := make([]O, 0)
		for i := 0; i < step.chunkSize; i++ {
			var item I
			var value O
			if item, err = step.reader(); err != nil {
				return
			}
			log.Println("read item: ", item)
			if reflect.ValueOf(item).IsZero() {
				//err = ErrStepNoValueToProcess
				break
			}
			if step.processor != nil {
				log.Println("process item: ", item)
				if value, err = step.processor(item); err != nil {
					return
				}
			} else {
				if data, ok := any(item).(O); ok {
					value = data
				}
			}
			log.Printf("add item %v to chunk: %v\n", item, chunckedItems)
			chunckedItems = append(chunckedItems, value)
		}
		if len(chunckedItems) == 0 {
			break
		}
		if err = step.writer(chunckedItems); err != nil {
			return
		}
	}
	log.Printf("Finishing step %v\n", step.name)
	return
} */

func (step *DefaultStep[I, O]) WithChunkSizeOf(size int) *DefaultStep[I, O] {
	step.chunkSize = size
	return step
}

func NewStep[I any, O any](reader ItemReader[I], processor ItemProcessor[I, O], writer ItemWriter[O]) *DefaultStep[I, O] {
	stepInstance := DefaultStep[I, O]{}
	stepInstance.name = uuid.NewString()
	stepInstance.reader = reader
	stepInstance.writer = writer
	stepInstance.chunkSize = DefaultStepChunkSize
	stepInstance.processor = processor
	return &stepInstance
}
func NewSimpleStep[I any](reader ItemReader[I], writer ItemWriter[I]) *DefaultStep[I, I] {
	stepInstance := DefaultStep[I, I]{}
	stepInstance.name = uuid.NewString()
	stepInstance.reader = reader
	stepInstance.writer = writer
	stepInstance.chunkSize = DefaultStepChunkSize
	return &stepInstance
}
