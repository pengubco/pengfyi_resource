package runners

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"peng.fyi/simple-loadtest/tasks"
)

// TaskRunner runs a task's Run() periodically. After each run, the runner waits for
// sometime, // determined by a interval with jitter. The runner can also registers a
// onStop() to run right before the runner stops.
type TaskRunner struct {
	id     uuid.UUID
	logger *zap.Logger

	lifetime      time.Duration
	cycleInterval time.Duration
	jitter        time.Duration

	// callback when the task runner stopped.
	onStop func()

	startOnce sync.Once
}

func NewTaskRunner(
	logger *zap.Logger,
	lifetime time.Duration,
	options ...TaskRunnerOption,
) *TaskRunner {
	r := TaskRunner{
		id:            uuid.New(),
		logger:        logger,
		lifetime:      lifetime,
		cycleInterval: 30 * time.Second,
		jitter:        5 * time.Second,
	}
	r.WithOptions(options...)
	return &r
}

// TaskRunnerOption adds optional behaviors to TaskRunner.
type TaskRunnerOption func(*TaskRunner)

// WithOnStop sets the callback when the runner stops.
func WithOnStop(onStop func()) TaskRunnerOption {
	return func(r *TaskRunner) {
		r.onStop = onStop
	}
}

// WithIntervalAndJitter sets the interval and jitter before next run.
func WithIntervalAndJitter(interval time.Duration, jitter time.Duration) TaskRunnerOption {
	return func(r *TaskRunner) {
		r.cycleInterval = interval
		r.jitter = jitter
	}
}

func (r *TaskRunner) WithOptions(options ...TaskRunnerOption) {
	for _, opt := range options {
		opt(r)
	}
}

// Start runs the task till ctx is canceled.
func (r *TaskRunner) Start(ctx context.Context, task tasks.Task) {
	l := r.logger.With(zap.String("runner-id", r.id.String()))

	periodicTicker := time.NewTicker(r.cycleInterval)
	stop := func() {
		periodicTicker.Stop()
		l.Info("stopping runner, cleanup")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		err := task.Cleanup(ctx)
		if err != nil && err != context.Canceled {
			l.Error("task cleanup error", zap.Error(err))
		}
		if r.onStop != nil {
			r.onStop()
		}
	}

	r.startOnce.Do(func() {
		ctx, cancel := context.WithTimeout(ctx, r.lifetime)
		defer cancel()
		if err := task.Setup(ctx); err != nil {
			l.Error("cannot setup task", zap.Error(err))
			return
		}
		l.Info("task setup success. going to run in cycles")
		for {
			select {
			case <-ctx.Done():
				stop()
				return
			case <-periodicTicker.C:
				l.Info("run task in a new cycle")
				err := task.Run(ctx)
				if err != nil && err != context.Canceled {
					l.Error("task run error", zap.Error(err))
				}
				periodicTicker.Reset(r.cycleInterval +
					time.Duration(rand.Int63n(int64(r.jitter))))
			}
		}
	})
}
