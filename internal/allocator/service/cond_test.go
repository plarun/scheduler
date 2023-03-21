package service

import (
	"testing"

	"github.com/plarun/scheduler/api/types/entity/task"
)

func TestConditionCheck(t *testing.T) {
	ts1 := task.NewTaskStatus(1, "task1", task.StateSuccess)
	ts2 := task.NewTaskStatus(2, "task2", task.StateSuccess)
	ts3 := task.NewTaskStatus(3, "task3", task.StateSuccess)
	ts4 := task.NewTaskStatus(3, "task3", task.StateFailure)
	ts5 := task.NewTaskStatus(4, "task4", task.StateFailure)
	ts6 := task.NewTaskStatus(4, "task4", task.StateSuccess)
	ts7 := task.NewTaskStatus(1, "task1", task.StateFailure)

	tests := map[string]struct {
		condition string
		stat      []*task.TaskStatus
		want      bool
	}{
		"good 1": {condition: "su(task1)", stat: []*task.TaskStatus{ts1}, want: true},
		"good 2": {condition: "su(task2)", stat: []*task.TaskStatus{ts2}, want: true},
		"good 3": {condition: "su(task3)", stat: []*task.TaskStatus{ts3}, want: true},
		"good 4": {
			condition: "su(task1)&su(task2)&(fa(task3)|nr(task4))",
			stat:      []*task.TaskStatus{ts1, ts2, ts4, ts5}, want: true},
		"good 5": {
			condition: "su(task1)&su(task2)&su(task3)&su(task4)",
			stat:      []*task.TaskStatus{ts1, ts2, ts3, ts6}, want: true},
		"good 6": {
			condition: "su(task1)&su(task2)&su(task3)&su(task4)",
			stat:      []*task.TaskStatus{ts1, ts2, ts3, ts5}, want: false},
		"good 7": {
			condition: "su(task1)&su(task2)&su(task3)&su(task4)",
			stat:      []*task.TaskStatus{ts7, ts2, ts3, ts6}, want: false},
		"good 8": {condition: "", stat: []*task.TaskStatus{}, want: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			chk := task.NewConditionChecker(0, tc.condition, tc.stat)

			if got, err := chk.Check(); err != nil {
				t.Fatalf("Check: %v", err)
			} else if got != tc.want {
				t.Fatalf("want: %#v, got: %#v", tc.want, got)
			}
		})
	}
}
