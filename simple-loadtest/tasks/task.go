package tasks

import "context"

// Task describes a task to be run by a TaskRunner.
// Each task has three stages, moving in one direction.
//  1. Setup: Get the task ready.
//  2. Run: TaskRunner calls Run periodically.
//  3. Cleanup.
//
// Each method must return when the context is canceled.
type Task interface {
	Setup(context.Context) error

	Run(context.Context) error

	Cleanup(context.Context) error
}
