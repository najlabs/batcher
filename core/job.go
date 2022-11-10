package core

import (
	"log"

	"github.com/google/uuid"
)

type Job interface {
	Run() error
}

type defaultJob struct {
	name  string
	steps []Step
}

func (job *defaultJob) Run() error {
	log.Printf("Start running job: %v with %d\n", job.name, len(job.steps))
	for _, step := range job.steps {
		log.Println("----------------------------")
		if err := step.Execute(); err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func NewJob(steps ...Step) Job {
	return &defaultJob{name: uuid.NewString(), steps: steps}
}
