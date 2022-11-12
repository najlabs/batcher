package batcher

type ExecutionContext struct {
	store map[string]interface{}
}
type JobExecutionContext struct {
	ExecutionContext
	CurrentStep string
}

func NewJobExecutionContext() *JobExecutionContext {
	return &JobExecutionContext{ExecutionContext{store: make(map[string]interface{})}, ""}
}
