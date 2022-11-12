package batcher

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"testing"
	"time"
)

func simpleStateReader() ItemReader[int] {
	index := 2
	return func(buffer []int) (int, error) {
		itemsRead := 0
		for i := 0; i < len(buffer) && index > 0; i++ {
			buffer[i] = index
			index--
			time.Sleep(time.Second * 1)
			itemsRead++
		}
		if index <= 0 {
			return itemsRead, io.EOF
		}
		return itemsRead, nil
	}
}

var simpleReader ItemReader[int] = func([]int) (int, error) { return 1, nil }
var simpleReaderErr ItemReader[int] = func([]int) (int, error) { return 0, errors.New("err1") }

// var simpleProcessor ItemMapper[int, int] = func(i int) (int, error) { return i + 1, nil }
// var simpleProcessor2 ItemMapper[int, string] = func(i int) (string, error) { return fmt.Sprintf("%d", i+1), nil }
var consoleWriter ItemWriter[int] = func(i []int) error { fmt.Printf("write: %v\n", i); return nil }

// var consoleWriter1 ItemWriter[string] = func(i []string) error { fmt.Printf("write: %v\n", i); return nil }
var consoleWriterErr ItemWriter[int] = func(i []int) error { fmt.Printf("write: %v\n", i); return errors.New("err") }

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
		{"simple step without option", args[int]{simpleReader, consoleWriter, nil}, &UniStep[int]{nil, simpleReader, consoleWriter}},
		{"simple step with option", args[int]{simpleReader, consoleWriter, []StepOption{}}, &UniStep[int]{nil, simpleReader, consoleWriter}},
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
		{"reader execution error", fields[int]{&coreStep{chunkSize: 1}, simpleReaderErr, consoleWriter}, true},
		{"writer execution error", fields[int]{&coreStep{chunkSize: 1}, simpleStateReader(), consoleWriterErr}, true},
		{"execution without error", fields[int]{&coreStep{chunkSize: 1}, simpleStateReader(), consoleWriter}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			step := &UniStep[int]{
				prop:   tt.fields.prop,
				reader: tt.fields.reader,
				writer: tt.fields.writer,
			}
			if err := step.execute(); (err != nil) != tt.wantErr {
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
		{"should have step1 and size 1", NewStep(simpleStateReader(), consoleWriter, WithName("step1"), WithChunkSize(1)), "step1", 1},
		{"should have step2 name and size 90", NewStep(simpleStateReader(), consoleWriter, WithName("step2"), WithChunkSize(90)), "step2", 90},
		{"should have step2 name and default size 90", NewStep(simpleStateReader(), consoleWriter, WithName("step2")), "step2", 10},
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
