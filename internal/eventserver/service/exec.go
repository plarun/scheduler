package service

import (
	"context"
	"fmt"

	"github.com/plarun/scheduler/internal/eventserver/db/query"
)

func PullReadyTasks(ctx context.Context) ([]int64, error) {
	var empty []int64
	// lock newly inserted ready tasks for pull
	if err := query.SwitchLockReadyTasks(0, 1); err != nil {
		return empty, fmt.Errorf("PullReadyTasks: %w", err)
	}

	// pull locked tasks for execution
	tasks, err := query.PullReadyTasks()
	if err != nil {
		return empty, fmt.Errorf("PullReadyTasks: %w", err)
	}

	// lock as tasks pulled for execution
	if err := query.SwitchLockReadyTasks(1, 2); err != nil {
		return empty, fmt.Errorf("PullReadyTasks: %w", err)
	}
	return tasks, nil
}

func GetTaskCommand(ctx context.Context, id int64) (string, string, string, error) {
	if command, fout, ferr, err := query.GetTaskCommand(id); err != nil {
		return "", "", "", fmt.Errorf("GetTaskCommand: %w", err)
	} else {
		return command, fout, ferr, nil
	}
}
