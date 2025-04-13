package interfaces

type SchedulerTask interface {
	ScheduledRun() error
}
