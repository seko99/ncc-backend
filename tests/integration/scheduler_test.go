package integration

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/logger/zero"
	scheduler2 "code.evixo.ru/ncc/ncc-backend/pkg/scheduler"
	"code.evixo.ru/ncc/ncc-backend/services/scheduler"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type TestTask struct {
	callCount int
}

func (s *TestTask) ScheduledRun() error {
	s.callCount++
	time.Sleep(3 * time.Second)
	return nil
}

func TestSchedulerSingleTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := zero.NewLogger()

	schedulerService := scheduler.NewScheduler(log)

	testTask := &TestTask{}

	schedulerService.RegisterTask(scheduler2.Task{
		Task:      testTask,
		Name:      "Test task",
		IsEnabled: true,
		Schedule: scheduler2.Schedule{
			Every: 2 * time.Second,
		},
	})

	err := schedulerService.Run(10 * time.Second)
	assert.NoError(t, err)
	assert.Equal(t, 5, testTask.callCount)
}
