package main

import (
	"fmt"
	"go-masters/task1/schedule"
	"time"
)

func main() {
	schedule.Add(&MyTask{}, time.Now())
	schedule.Add(&MyTask{}, time.Now().Add(time.Second))
	schedule.RunSchedule()
}

type MyTask struct {
}

func (m *MyTask) Exec() {
	fmt.Println(time.Now())
}
