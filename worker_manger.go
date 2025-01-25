package workermanager

import (
	"context"
	"fmt"
	"sync"
)

type WorkManager interface {
	AddWorker(worker Worker) error
	Start(name string) error
	StartAll() error
	Stop(name string) (chan struct{}, error)
	StopAll() chan struct{}
}

type workerManager struct {
	ctx     context.Context
	workers map[string]Worker
	mu      sync.RWMutex
}

func NewWorkerManager(ctx context.Context) WorkManager {
	return &workerManager{
		ctx:     ctx,
		workers: make(map[string]Worker),
		mu:      sync.RWMutex{},
	}
}

func (wm *workerManager) AddWorker(worker Worker) error {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	if _, ok := wm.workers[worker.Name()]; ok {
		return fmt.Errorf("worker %s already exists", worker.Name())
	}

	wm.workers[worker.Name()] = worker

	return nil
}

func (wm *workerManager) Start(name string) error {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	worker, ok := wm.workers[name]
	if !ok {
		return fmt.Errorf("worker %s not found", name)
	}

	return worker.Start(wm.ctx)
}

func (wm *workerManager) StartAll() error {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	for _, worker := range wm.workers {
		worker.Start(wm.ctx)
	}

	return nil
}

func (wm *workerManager) Stop(name string) (chan struct{}, error) {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	worker, ok := wm.workers[name]
	if !ok {
		return nil, fmt.Errorf("worker %s not found", name)
	}

	return worker.Stop(context.Background()), nil
}

func (wm *workerManager) StopAll() chan struct{} {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	stop := make(chan struct{})
	defer close(stop)

	wg := sync.WaitGroup{}
	for _, worker := range wm.workers {
		wg.Add(1)
		go func(worker Worker) {
			<-worker.Stop(context.Background())
			wg.Done()
		}(worker)
	}
	wg.Wait()

	return stop
}
