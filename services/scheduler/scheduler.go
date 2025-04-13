package scheduler

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/pkg/scheduler"
	"sync"
	"time"
)

type Scheduler struct {
	sync.RWMutex
	log   logger.Logger
	tasks map[string]scheduler.Task
}

func (s *Scheduler) RegisterTask(task scheduler.Task) {
	s.tasks[task.Name] = task
}

func (s *Scheduler) Run(maxDuration ...time.Duration) error {
	startTime := time.Now()
	for {
		for k, task := range s.tasks {
			if task.IsEnabled {
				now := time.Now()
				nextRun := task.LastRun.Add(task.Schedule.Every)

				if task.Schedule.Every > 0 && now.After(nextRun) {
					task.LastRun = now
					s.Lock()
					s.tasks[k] = task
					s.Unlock()

					cTask := task
					cK := k
					go func() {
						err := cTask.Task.ScheduledRun()
						if err != nil {
							s.log.Error("Error executing task %s: %v", cK, err)
						}
					}()
				}
			}
		}
		time.Sleep(time.Second)
		if len(maxDuration) > 0 {
			elapsed := time.Now().UnixNano() - startTime.UnixNano()
			if elapsed >= maxDuration[0].Nanoseconds() {
				return nil
			}
		}
	}
}

func NewScheduler(
	log logger.Logger,
) *Scheduler {
	return &Scheduler{
		log:   log,
		tasks: map[string]scheduler.Task{},
	}
}
