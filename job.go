package batcher

import (
	"log"
)

type JobConfig func(*defaultJob)
type JobListener func(*JobExecutionContext)
type defaultJob struct {
	name           string
	steps          []Step
	beforeAllStep  JobListener
	afterAllStep   JobListener
	beforeEachStep JobListener
	afterEachStep  JobListener
}

func (job *defaultJob) Steps(steps ...Step) {
	job.steps = steps
}

func (job *defaultJob) Configure(configs ...JobConfig) {
	for _, config := range configs {
		config(job)
	}
}

func (job *defaultJob) Run() {
	log.Printf("Start running job: %v with %d steps\n", job.name, len(job.steps))
	ctx := NewJobExecutionContext()
	if job.beforeAllStep != nil {
		job.beforeAllStep(ctx)
	}
	for _, step := range job.steps {
		log.Println("----------------------------")
		ctx.CurrentStep = step.name()
		if job.beforeEachStep != nil {
			job.beforeEachStep(ctx)
		}
		if err := step.execute(); err != nil {
			log.Fatal(err)
		}
		if job.afterEachStep != nil {
			job.afterEachStep(ctx)
		}
	}
	if job.afterAllStep != nil {
		job.afterAllStep(ctx)
	}
}

func NewJobInline(name string, steps ...Step) *defaultJob {
	return &defaultJob{name: name, steps: steps}
}
func NewJob(name string, configs ...JobConfig) *defaultJob {
	job := &defaultJob{name: name}
	for _, config := range configs {
		config(job)
	}
	return job
}

func BeforeEachStep(beforeEachListener JobListener) JobConfig {
	return func(job *defaultJob) {
		job.beforeEachStep = beforeEachListener
	}
}
func BeforeAllStep(beforeAllListener JobListener) JobConfig {
	return func(job *defaultJob) {
		job.beforeAllStep = beforeAllListener
	}
}
func AfterAllStep(afterAllListener JobListener) JobConfig {
	return func(job *defaultJob) {
		job.afterAllStep = afterAllListener
	}
}
func AfterEachStep(afterEachListener JobListener) JobConfig {
	return func(job *defaultJob) {
		job.afterEachStep = afterEachListener
	}
}
