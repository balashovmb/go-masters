package schedule

import (
	"slices"
	"sync"
	"time"
)

type Task interface {
	Exec()
}

type Job struct {
	task Task
	time time.Time
}

type Schedule struct {
	mu      sync.Mutex
	jobs    []Job
	running bool
}

var s = &Schedule{running: false}

func Add(task Task, time time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs = append(s.jobs, Job{task, time})
}

func RunSchedule() {
	if s.running {
		return
	}
	s.running = true
	for {
		if len(s.jobs) == 0 {
			s.running = false
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
