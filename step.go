package batcher

import (
	"errors"
	"github.com/lithammer/shortuuid/v4"
	"log"
)

const defaultChunkSize = 10

var ErrEnd = errors.New("ErrEnd")

// ItemReader read data from custom source and put in the buffer
// return the number of data readed and error if any.
// returning error batcher.ErrEnd will mean you ar done with reading
// type ItemReader[INPUT any] func([]INPUT) (int, error)
type ItemReader[INPUT any] interface {
	Read([]INPUT) (int, error)
}
type ItemWriter[OUTPUT any] interface {
	Write([]OUTPUT) error
}
type Step[IOType any] interface {
	Read([]IOType) (int, error)
	Write([]IOType) error
}

// type ItemWriter[OUTPUT any] func([]OUTPUT) error
type ItemMapper[INPUT any, OUTPUT any] func(INPUT) (OUTPUT, error)
type StepOption func(*coreStep)
type ExecutionListener func(*ExecutionContext)
type StepExecutor interface {
	execute() error
	name() string
}
type simpleStep[INPUT any, OUTPUT any] struct {
	prop   *coreStep
	reader ItemReader[INPUT]
	writer ItemWriter[OUTPUT]
	mapper ItemMapper[INPUT, OUTPUT]
}
type coreStep struct {
	name         string
	chunkSize    int
	beforeReader ExecutionListener
	afterReader  ExecutionListener
	beforeWriter ExecutionListener
	afterWriter  ExecutionListener
}

type UniStep[IOType any] struct {
	prop   *coreStep
	reader ItemReader[IOType]
	writer ItemWriter[IOType]
}

func (step *UniStep[IOType]) name() string {
	return step.prop.name
}
func (step *UniStep[IOType]) execute() error {
	log.Printf("Executing step %v\n", step.prop.name)
	for {
		chunckedItems := make([]IOType, step.prop.chunkSize)
		n, err := step.reader.Read(chunckedItems)
		if n > 0 {
			if werr := step.writer.Write(chunckedItems[:n]); werr != nil {
				return werr
			}
		}
		if err == ErrEnd {
			break
		}
		if err != nil {
			return err
		}
	}
	log.Printf("Finishing step %v\n", step.prop.name)
	return nil
}

func NewStepAction[IOType any](action Step[IOType], options ...StepOption) *UniStep[IOType] {
	return NewStep[IOType](action, action, options...)
}
func NewStep[IOType any](reader ItemReader[IOType], writer ItemWriter[IOType], options ...StepOption) *UniStep[IOType] {
	var uniStep = UniStep[IOType]{}
	uniStep.reader = reader
	uniStep.writer = writer
	uniStep.prop = &coreStep{}
	uniStep.prop.name = shortuuid.New()
	uniStep.prop.chunkSize = defaultChunkSize
	for _, option := range options {
		option(uniStep.prop)
	}
	return &uniStep
}

func WithName(name string) StepOption {
	return func(c *coreStep) {
		c.name = name
	}
}
func WithChunkSize(size int) StepOption {
	return func(c *coreStep) {
		c.chunkSize = size
	}
}
