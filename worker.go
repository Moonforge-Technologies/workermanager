package workermanager

import (
	"context"
)

type Status string

const (
	StatusRunning Status = "running"
	StatusStopped Status = "stopped"
)

type Worker interface {
	Name() string
	Status() Status
	Start(ctx context.Context) error
	Stop(ctx context.Context) chan struct{}
}
