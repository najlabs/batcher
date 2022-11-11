package core

import (
	"fmt"
)

const StepHeader = "Step"

var errDone = newStepState("current execution done")
var errContinue = newStepState("continue current execution")
var errCommitBuffered = newStepState("commit buffered items")
var errSkip = newStepState("skip current execution")
var errCancel = newStepState("cancel current execution")

type State struct {
	header  string
	message string
}

func (err State) Error() string {
	if err.header == "" {
		return err.Error()
	}
	return fmt.Sprintf("[%v] %v", err.header, err.message)
}

/*
	func newStateWithError(err error) State {
		return State{"", err.Error()}
	}
*/
/*func newState(header string, message string) State {
	return State{header, message}
}*/

func newStepState(message string) State {
	return State{StepHeader, message}
}
