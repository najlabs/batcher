package core

import (
	"log"
)

type BatchContext struct {
	store map[string]interface{}
}

func (ctx *BatchContext) UpdateIfExists(key string, value interface{}) {
	_, ok := ctx.store[key]
	if ok {
		ctx.store[key] = value
	}
}

func (ctx *BatchContext) PutIfNotExists(key string, value interface{}) {
	_, ok := ctx.store[key]
	if !ok {
		ctx.store[key] = value
	}
}
func (ctx *BatchContext) GetOrCreate(key string, defaultValue interface{}) interface{} {
	value, ok := ctx.store[key]
	if ok {
		return value
	}
	ctx.store[key] = defaultValue
	return defaultValue
}
func (ctx *BatchContext) GetOrElse(key string, defaultValue interface{}) interface{} {
	value, ok := ctx.store[key]
	if !ok {
		return defaultValue
	}
	return value
}
func (ctx *BatchContext) nextState() State {
	return errContinue
}
func (ctx *BatchContext) Done() State {
	return errDone
}

/*
	func (ctx *BatchContext) Cancel() State {
		return errCancel
	}
*/
func (ctx *BatchContext) Skip() State {
	return errSkip
}

type ReaderContext[I any] struct {
	BatchContext
	bufferedItems []I
	bufferSize    int
}

func (ctx *ReaderContext[I]) ResetBuffer() {
	ctx.bufferedItems = make([]I, 0, ctx.bufferSize)
}
func (ctx *ReaderContext[I]) commitState() State {
	if len(ctx.bufferedItems) == 0 {
		return errDone
	}
	return errCommitBuffered
}
func (ctx *ReaderContext[I]) Done() State {
	if len(ctx.bufferedItems) > 0 {
		return ctx.commitState()
	}
	return errDone
}
func (ctx *ReaderContext[I]) Process(item I) State {
	ctx.bufferedItems = append(ctx.bufferedItems, item)
	log.Printf("Add Next item %v size %v buffer %v", item, len(ctx.bufferedItems), ctx.bufferedItems)
	if len(ctx.bufferedItems) == ctx.bufferSize {
		return ctx.commitState()
	}
	return ctx.nextState()
}

/*func NewContext(bufferSize int) *BatchContext {
	ctx := &BatchContext{}
	ctx.store = make(map[string]interface{})
	return ctx
}*/

func NewReaderContext[I any](bufferSize int) *ReaderContext[I] {
	ctx := &ReaderContext[I]{}
	ctx.store = make(map[string]interface{})
	ctx.bufferSize = bufferSize
	log.Printf("With buffersize %v\n", bufferSize)
	ctx.bufferedItems = make([]I, 0, bufferSize)
	return ctx
}

type WriterContext[O any] struct {
	BatchContext
	itemToWrite []O
}

func (ctx *WriterContext[O]) ItemsToWrite() []O {
	return ctx.itemToWrite
}

func NewWriterContext[O any](store map[string]interface{}) *WriterContext[O] {
	ctx := &WriterContext[O]{}
	ctx.store = store
	return ctx
}

type ProcessorContext[I any, O any] struct {
	BatchContext
	ReaderContext[I]
	ItemToWrite O
}
