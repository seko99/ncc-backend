package scheduler

import (
	"code.evixo.ru/ncc/ncc-backend/services/interfaces"
	"time"
)

type Task struct {
	Task      interfaces.SchedulerTask
	Name      string
	IsEnabled bool
	LastRun   time.Time
	Schedule  Schedule
}

type Schedule struct {
	Every time.Duration
}
