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

type ItemReader[I any] func() (I, error)
type ItemWriter[O any] func(o []O) error
type ItemProcessor[I any, O any] func(i I) (O, error)

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

func (step *DefaultStep[I, O]) Execute() (err error) {
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
}

func (step *DefaultStep[I, O]) WithChunkSizeOf(size int) *DefaultStep[I, O] {
	step.chunkSize = size
	return step
}

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
