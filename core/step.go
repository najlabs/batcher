package core

import (
	"errors"
	"log"
	"reflect"

	"github.com/google/uuid"
)

type context interface {
	Get(key string) interface{}
	Put(key string, value interface{})
}
type ItemReaderContext interface {
	context
	ReadItem()
}

var ErrStepNoValueToProcess error = errors.New("Reading has no more value")
var ErrMissingProcessorOrMapper error = errors.New("processor or mapper is required")

const DEFAULT_STEP_CHUNK_SIZE = 10

type Item map[string]any

//	type ItemReader[T any] interface {
//		Read() chan T
//	}
//
//	type ItemWriter[T any] interface {
//		Write(chan Item)
//	}
//
//	type ItemProcessor[I any, O any] interface {
//		Write(chan I) chan O
//	}
type ItemReader[I any] func() (I, error)
type ItemWriter[O any] func(o O) error
type ItemProcessor[I any, O any] func(i I) (O, error)
type ItemMapper[I any, O any] func(i I) (O, error)

// https://docs.spring.io/spring-batch/docs/current/reference/html/domain.html#domainLanguageOfBatch
type Step interface {
	Execute() error
}
type DefaultStep[I any, O any] struct {
	name      string
	reader    ItemReader[I]
	writer    ItemWriter[O]
	processor ItemProcessor[I, O]
	mapper    ItemMapper[I, O]
	chunkSize int
}

func (step *DefaultStep[I, O]) Execute() (err error) {
	log.Printf("Executing step %v\n", step.name)
	// for i := 0; i < DEFAULT_STEP_CHUNK_SIZE; i++ {
	for {
		var item I
		var value O
		if item, err = step.reader(); err != nil {
			return
		}
		if reflect.ValueOf(item).IsZero() {
			//err = ErrStepNoValueToProcess
			break
		}
		if step.processor != nil {

			if value, err = step.processor(item); err != nil {
				return
			}
		} else {
			if data, ok := any(item).(O); ok {
				value = data
			}

			/* if step.mapper == nil {
				err = ErrMissingProcessorOrMapper
				return
			}
			if value, err = step.mapper(item); err != nil {
				if err = step.writer(value); err != nil {
					return
				}
			} */
		}
		if err = step.writer(value); err != nil {
			return
		}
	}
	log.Printf("Finishing step %v\n", step.name)
	return
}

/* func (step *DefaultStep[I, O]) WithProcessor(writer ItemProcessor[I, O]) *DefaultStep[I, O] {

	return step
}

func (step *DefaultStep[I, O]) WithWriter(writer ItemWriter[O]) *DefaultStep[I, O] {
	return step
}

func (step *DefaultStep[I, O]) WithReader(reader ItemReader[I]) *DefaultStep[I, O] {
	return step
}

func NewStepBuilder[I any, O any]() *DefaultStep[I, O] {
	return &DefaultStep[I, O]{}
}
*/

func NewStep[I any, O any](reader ItemReader[I], processor ItemProcessor[I, O], writer ItemWriter[O]) *DefaultStep[I, O] {
	stepInstance := DefaultStep[I, O]{}
	stepInstance.name = uuid.NewString()
	stepInstance.reader = reader
	stepInstance.writer = writer
	stepInstance.chunkSize = DEFAULT_STEP_CHUNK_SIZE
	stepInstance.processor = processor
	return &stepInstance
}
func NewSimpleStep[I any](reader ItemReader[I], writer ItemWriter[I]) *DefaultStep[I, I] {
	stepInstance := DefaultStep[I, I]{}
	stepInstance.name = uuid.NewString()
	stepInstance.reader = reader
	stepInstance.writer = writer
	stepInstance.chunkSize = DEFAULT_STEP_CHUNK_SIZE
	return &stepInstance
}
