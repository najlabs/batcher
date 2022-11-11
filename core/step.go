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
func (step *DefaultStep[I, O]) executeStateReader(readerCtx *ReaderContext[I]) error {
	switch err := step.reader(readerCtx); err {
	case errContinue:
		return step.executeStateReader(readerCtx)
	case errCommitBuffered:
		log.Println("committed reader ")
		//err = step.commitBuffer(readerCtx)
		return nil
	case errSkip:
		log.Println("skip reader ")
		return step.executeStateReader(readerCtx)
	default:
		return err
	}
}
func (step *DefaultStep[I, O]) executeStateWriter(ctx *WriterContext[O]) error {
	switch err := step.writer(ctx); err {
	case errContinue:
		return step.executeStateWriter(ctx)
	case errDone:
		log.Println("Done writer ")
		//_ = step.commitBuffer(readerCtx)
		return nil
	case errSkip:
		log.Println("skip writer ")
		return step.executeStateWriter(ctx)
	default:
		return err
	}
}
func (step *DefaultStep[I, O]) Execute() error {
	log.Printf("Executing step %v\n", step.name)
	readerCtx := NewReaderContext[I](step.chunkSize)
	for {
		if err := step.executeStateReader(readerCtx); err != nil {
			if errors.Is(err, errDone) {
				break
			}
			return err
		}
		writerCtx := NewWriterContext[O](readerCtx.store)
		var err error
		if writerCtx.itemToWrite, err = step.mapIO(readerCtx.bufferedItems); err != nil {
			return err
		}
		if err = step.executeStateWriter(writerCtx); err != nil {
			return err
		}
		readerCtx.ResetBuffer()
	}
	log.Printf("Finishing step %v\n", step.name)
	return nil
}

func (step *DefaultStep[I, O]) WithChunkSizeOf(size int) *DefaultStep[I, O] {
	step.chunkSize = size
	return step
}

/*
	func NewStep[I any, O any](reader ItemReader[I], processor ItemProcessor[I, O], writer ItemWriter[O]) *DefaultStep[I, O] {
		stepInstance := DefaultStep[I, O]{}
		stepInstance.name = uuid.NewString()
		stepInstance.reader = reader
		stepInstance.writer = writer
		stepInstance.chunkSize = DefaultStepChunkSize
		stepInstance.processor = processor
		return &stepInstance
	}
*/
func NewSimpleStep[I any](reader ItemReader[I], writer ItemWriter[I]) *DefaultStep[I, I] {
	stepInstance := DefaultStep[I, I]{}
	stepInstance.name = uuid.NewString()
	stepInstance.reader = reader
	stepInstance.writer = writer
	stepInstance.chunkSize = DefaultStepChunkSize
	return &stepInstance
}
