package core

import (
	"log"

	"github.com/google/uuid"
)

type Job interface {
	Run()
}

type defaultJob struct {
	name  string
	steps []Step
}

func (job *defaultJob) Run() {
	log.Printf("Start running job: %v with %d steps\n", job.name, len(job.steps))
	for _, step := range job.steps {
		log.Println("----------------------------")
		if err := step.Execute(); err != nil {
			panic(err)
		}
	}
}

func NewJob(steps ...Step) Job {
	return &defaultJob{name: uuid.NewString(), steps: steps}
}
