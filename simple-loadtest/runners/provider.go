package runners

import (
	"github.com/samber/lo"
	"peng.fyi/simple-loadtest/tasks"
)

// TaskRunnerProvider provides factory functions that create Task and TaskRunner.
type TaskRunnerProvider interface {
	NextTask() func() (tasks.Task, error)

	NextRunner() func() (*TaskRunner, error)
}

// TaskRunnerRR povides Task and TaskRunner in pairs, in round robin manner.
type TaskRunnerRR struct {
	newRunners *RoundRobinContainer[func() (*TaskRunner, error)]
	newTasks   *RoundRobinContainer[func() (tasks.Task, error)]
}

func NewTaskRunnerRR(newRunners []func() (*TaskRunner, error), newTasks []func() (tasks.Task, error)) *TaskRunnerRR {
	return &TaskRunnerRR{
		newRunners: lo.Must(NewRoundRobinContainer(newRunners...)),
		newTasks:   lo.Must(NewRoundRobinContainer(newTasks...)),
	}
}

func (p *TaskRunnerRR) NextRunner() func() (*TaskRunner, error) {
	return p.newRunners.Next()
}

func (p *TaskRunnerRR) NextTask() func() (tasks.Task, error) {
	return p.newTasks.Next()
}
