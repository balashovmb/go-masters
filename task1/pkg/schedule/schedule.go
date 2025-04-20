package schedule

import (
	"slices"
	"sync"
	"time"
)

type Taskable interface {
	Exec()
}

type Job struct {
	task Taskable
	time time.Time
}

type Schedule struct {
	mu   sync.Mutex
	jobs []Job
}

var s = &Schedule{}

func Add(task Taskable, time time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs = append(s.jobs, Job{task, time})
}

func RunSchedule() {
	for {
		if len(s.jobs) == 0 {
			return
		}

		for _, job := range s.jobs {
			if time.Now().After(job.time) {
				go runJob(job)
			}
		}

		time.Sleep(1000 * time.Millisecond)
	}
}

func runJob(job Job) {
	job.task.Exec()
	removeJob(job)
}

func removeJob(job Job) {
	s.mu.Lock()
	defer s.mu.Unlock()
	index := slices.Index(s.jobs, job)
	s.jobs = slices.Delete(s.jobs, index, index+1)
}
