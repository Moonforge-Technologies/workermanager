package workermanager_test

import (
	"context"
	"testing"
	"time"
	"workermanager"
	"workermanager/internal"

	"github.com/stretchr/testify/suite"
)

type WorkerManagerTestSuite struct {
	suite.Suite
}

func TestWorkerManager(t *testing.T) {
	suite.Run(t, new(WorkerManagerTestSuite))
}

func (s *WorkerManagerTestSuite) TestStopAll_ShouldWaitAllWorkersToStop() {
	// Arrange
	var workers = []workermanager.Worker{
		internal.NewWorker("worker-1"),
		internal.NewWorker("worker-2"),
		internal.NewWorker("worker-3"),
		internal.NewWorker("worker-4"),
		internal.NewWorker("worker-5"),
	}
	var wm = workermanager.NewWorkerManager(context.Background())
	for _, worker := range workers {
		wm.AddWorker(worker)
		wm.Start(worker.Name())
	}
	time.Sleep(1 * time.Second)

	// Act
	<-wm.StopAll()

	// Assert
	for _, worker := range workers {
		s.Equal(workermanager.StatusStopped, worker.Status())
	}
}

func (s *WorkerManagerTestSuite) TestWorkerManager_WhenStop_ShouldAwaitWorkerStop() {
	// Arrange
	var worker = internal.NewWorker("worker-1")
	var wm = workermanager.NewWorkerManager(context.Background())
	wm.AddWorker(worker)
	wm.Start(worker.Name())
	time.Sleep(1 * time.Second)

	// Act
	var stopChan, err = wm.Stop(worker.Name())
	s.Require().NoError(err)
	<-stopChan

	// Assert
	s.Equal(workermanager.StatusStopped, worker.Status())
}
