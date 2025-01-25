package internal

import (
	"context"
	"fmt"
	"sync"
	"time"
	"workermanager"
)

type workerMock struct {
	name      string
	status    workermanager.Status
	ctx       context.Context
	ctxCancel context.CancelFunc
	stopChan  chan struct{}
	mu        sync.Mutex
}

func NewWorker(name string) workermanager.Worker {
	return &workerMock{
		name:     name,
		status:   workermanager.StatusStopped,
		stopChan: make(chan struct{}),
		mu:       sync.Mutex{},
	}
}

func (w *workerMock) Name() string {
	return w.name
}

func (w *workerMock) Status() workermanager.Status {
	return w.status
}

func (w *workerMock) Start(ctx context.Context) error {
	w.ctx, w.ctxCancel = context.WithCancel(ctx)

	go w.goRoutine()

	return nil
}

func (w *workerMock) Stop(ctx context.Context) chan struct{} {
	w.ctxCancel()
	return w.stopChan
}

func (w *workerMock) goRoutine() {
	defer func() {
		fmt.Println("worker", w.Name(), "is stopping")
		w.setStatus(workermanager.StatusStopped)
		close(w.stopChan)
	}()

	w.setStatus(workermanager.StatusRunning)

	for {
		select {
		case <-w.ctx.Done():
			fmt.Println("worker", w.Name(), "context is done")
			return
		default:
			fmt.Println("worker", w.Name(), " is working")
			time.Sleep(5 * time.Second)
		}
	}
}

func (w *workerMock) setStatus(status workermanager.Status) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.status = status
}
