package service

import (
	"context"
	"fmt"

	"github.com/plarun/scheduler/api/types/entity/task"
	"github.com/plarun/scheduler/internal/eventserver/db/query"
)

func ChangeTaskState(ctx context.Context, id int64, state task.State) error {
	// change the task state
	if err := query.SetTaskStatus(id, state); err != nil {
		return fmt.Errorf("ChangeTaskState: %w", err)
	}

	// remove from ready queue
	if state.IsRunning() {
		if err := query.RemoveFromReady(id); err != nil {
			return fmt.Errorf("ChangeTaskState: %w", err)
		}
	}

	// unstage the task
	if state.IsFailure() || state.IsSuccess() || state.IsAborted() || state.IsFrozen() {
		if err := query.UnstageTask(id); err != nil {
			return fmt.Errorf("ChangeTaskState: %w", err)
		}

		// unlock the task
		if err := query.UnlockUnstagedTask(id); err != nil {
			return fmt.Errorf("ChangeTaskState: %w", err)
		}

		// if its bundle task is in staging area then
		// check if its any of tasks are running else
		// set the status for bundle and unstage it

		hasParent, parentId, badRuns, err := query.BundleAndSiblingsStatus(id)
		if err != nil {
			return fmt.Errorf("ChangeTaskState: %w", err)
		}
		if !hasParent {
			return nil
		}

		pend, err := query.HasStagedSiblings(id)
		if err != nil {
			return fmt.Errorf("ChangeTaskState: %w", err)
		}

		if !state.IsFrozen() {
			if state.IsSuccess() {
				if !pend && badRuns == 0 {
					if err := query.SetTaskStatus(parentId, task.StateSuccess); err != nil {
						return fmt.Errorf("ChangeTaskState: %w", err)
					}
				}
			} else if state.IsAborted() || state.IsFailure() || badRuns > 0 {
				if err := query.SetTaskStatus(parentId, task.StateSuccess); err != nil {
					return fmt.Errorf("ChangeTaskState: %w", err)
				}
			}

			// unstage bundle
			if err := query.UnstageTask(parentId); err != nil {
				return fmt.Errorf("ChangeTaskState: %w", err)
			}

			// unlock bundle
			if err := query.UnlockUnstagedTask(parentId); err != nil {
				return fmt.Errorf("ChangeTaskState: %w", err)
			}
		}
	}
	return nil
}
