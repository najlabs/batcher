package batcher

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"
)

type dummyStateReader struct {
	index int
}

func (d *dummyStateReader) Read(buffer []int) (int, error) {
	itemsRead := 0
	for i := 0; i < len(buffer) && d.index > 0; i++ {
		buffer[i] = d.index
		d.index--
		time.Sleep(time.Second * 1)
		itemsRead++
	}
	if d.index <= 0 {
		return itemsRead, ErrEnd
	}
	return itemsRead, nil
}
func (d *dummyStateReader) Write(data []int) error {
	fmt.Printf("write: %v\n", data)
	return nil
}

type errReaderAndWriter struct{}

func (d errReaderAndWriter) Read(buffer []int) (int, error) {
	log.Println(buffer)
	return 0, errors.New("err1")
}
func (d errReaderAndWriter) Write(data []int) error {
	fmt.Printf("write: %v\n", data)
	return errors.New("err")
}

var simpleReaderWriter = &dummyStateReader{2}
var errorReaderWriter = &errReaderAndWriter{}

func TestNewStep(t *testing.T) {
	type args[IType any] struct {
		reader  ItemReader[IType]
		writer  ItemWriter[IType]
		options []StepOption
	}
	tests := []struct {
		name string
		args args[int]
		want *UniStep[int]
	}{
		{"simple step without option", args[int]{simpleReaderWriter, simpleReaderWriter, nil}, &UniStep[int]{nil, simpleReaderWriter, simpleReaderWriter}},
		{"simple step with option", args[int]{simpleReaderWriter, simpleReaderWriter, []StepOption{}}, &UniStep[int]{nil, simpleReaderWriter, simpleReaderWriter}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewStep(tt.args.reader, tt.args.writer, tt.args.options...)
			if reflect.ValueOf(got.prop).IsZero() {
				t.Errorf("prop is nil = %v", got.prop)
			}
			if reflect.ValueOf(got.reader).IsZero() {
				t.Errorf("NewStep() reader is nil= %v", got.reader)
			}
			if reflect.ValueOf(got.writer).IsZero() {
				t.Errorf("NewStep() writer is nil= %v", got.writer)
			}
		})
	}
}

func TestUniStep_execute(t *testing.T) {
	type fields[IType any] struct {
		prop   *coreStep
		reader ItemReader[IType]
		writer ItemWriter[IType]
	}
	tests := []struct {
		name    string
		fields  fields[int]
		wantErr bool
	}{
		{"reader execution error", fields[int]{&coreStep{chunkSize: 1}, errorReaderWriter, simpleReaderWriter}, true},
		{"writer execution error", fields[int]{&coreStep{chunkSize: 1}, simpleReaderWriter, errorReaderWriter}, true},
		{"execution without error", fields[int]{&coreStep{chunkSize: 1}, simpleReaderWriter, simpleReaderWriter}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			step := &UniStep[int]{
				prop:   tt.fields.prop,
				reader: tt.fields.reader,
				writer: tt.fields.writer,
			}
			if err := step.execute(); ((err != nil) != tt.wantErr) || ((err == ErrEnd) && !tt.wantErr) {
				t.Errorf("execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewStep_withNameOption(t *testing.T) {
	tests := []struct {
		name       string
		step       *UniStep[int]
		wantedName string
		wantedSize int
	}{
		{"should have step1 and size 1", NewStep[int](simpleReaderWriter, simpleReaderWriter, WithName("step1"), WithChunkSize(1)), "step1", 1},
		{"should have step2 name and size 90", NewStep[int](simpleReaderWriter, simpleReaderWriter, WithName("step2"), WithChunkSize(90)), "step2", 90},
		{"should have step2 name and default size 90", NewStep[int](simpleReaderWriter, simpleReaderWriter, WithName("step2")), "step2", 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.step.prop.name != tt.wantedName {
				t.Errorf("NewStep with incorrect name Option actual = %v, want %v", tt.step.prop.name, tt.wantedName)
			}
			if tt.step.prop.chunkSize != tt.wantedSize {
				t.Errorf("NewStep with incorrect name Option actual = %v, want %v", tt.step.prop.chunkSize, tt.wantedSize)
			}
		})
	}
}
func TestNewStepAction_withNameOption(t *testing.T) {
	tests := []struct {
		name       string
		step       *UniStep[int]
		wantedName string
		wantedSize int
	}{
		{"should have step1 and size 1", NewStepAction[int](simpleReaderWriter, WithName("step1"), WithChunkSize(1)), "step1", 1},
		{"should have step2 name and size 90", NewStepAction[int](simpleReaderWriter, WithName("step2"), WithChunkSize(90)), "step2", 90},
		{"should have step2 name and default size 90", NewStepAction[int](simpleReaderWriter, WithName("step2")), "step2", 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.step.prop.name != tt.wantedName {
				t.Errorf("NewStep with incorrect name Option actual = %v, want %v", tt.step.prop.name, tt.wantedName)
			}
			if tt.step.prop.chunkSize != tt.wantedSize {
				t.Errorf("NewStep with incorrect name Option actual = %v, want %v", tt.step.prop.chunkSize, tt.wantedSize)
			}
		})
	}
}
