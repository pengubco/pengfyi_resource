package runners

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// LoadRunner runs multiple tasks to generate load. It maintains a certain population
// of TaskRunners.
type LoadRunner struct {
	id uuid.UUID

	logger *zap.Logger

	lifetime time.Duration

	// Number of targeted active runners.
	targetedRunnerCnt int

	// Cancel func for each task runner. Right now, it's not been used because
	// the load runner cancel its own context, which is the parent of task runner's
	// context. The cancel func is here to support cancel specific task runner
	// before the task runner's context is canceled.
	runnerCancels map[uuid.UUID]func()

	runnerProvider TaskRunnerProvider

	stoppedRunnerCh chan uuid.UUID

	wg        sync.WaitGroup
	startOnce sync.Once
}

func NewLoadRunner(
	logger *zap.Logger,
	runnerProvider TaskRunnerProvider,
	lifetime time.Duration,
	targetedRunnerCnt int,
) *LoadRunner {
	r := LoadRunner{
		id:                uuid.New(),
		runnerProvider:    runnerProvider,
		lifetime:          lifetime,
		targetedRunnerCnt: targetedRunnerCnt,
		runnerCancels:     make(map[uuid.UUID]func()),
		stoppedRunnerCh:   make(chan uuid.UUID, targetedRunnerCnt),
	}
	r.logger = logger.With(zap.String("load-runner-id", r.id.String()))
	return &r
}

// Start maintains a population of active task runners.
func (r *LoadRunner) Start(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, r.lifetime)
	defer cancel()
	addRunner := func() {
		newRunner, newTask := r.runnerProvider.NextRunner(), r.runnerProvider.NextTask()
		runner, err := newRunner()
		if err != nil {
			r.logger.Error("cannot create task runner", zap.Error(err))
			return
		}
		task, err := newTask()
		if err != nil {
			r.logger.Error("cannot create task runner", zap.Error(err))
			return
		}
		r.logger.Info("created runner and task", zap.String("runner-id", runner.id.String()))

		runner.WithOptions(WithOnStop(func() {
			r.stoppedRunnerCh <- runner.id
		}))
		r.wg.Add(1)
		runnerCtx, runnerCancel := context.WithCancel(ctx)
		go func() {
			defer r.wg.Done()
			runner.Start(runnerCtx, task)
		}()
		r.runnerCancels[runner.id] = runnerCancel
	}

	r.startOnce.Do(func() {
		for len(r.runnerCancels) < r.targetedRunnerCnt {
			addRunner()
		}
		for {
			select {
			case <-ctx.Done():
				r.logger.Info("stopping load runner. wait for task runners to stop")
				r.wg.Wait()
				return
			case runnerID := <-r.stoppedRunnerCh:
				r.logger.Info("task runner stopped", zap.String("runner-id", runnerID.String()))
				delete(r.runnerCancels, runnerID)
				r.logger.Info("going to add a new runner")
				addRunner()
			}
		}
	})
}
