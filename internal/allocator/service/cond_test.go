package service

import (
	"testing"

	"github.com/plarun/scheduler/api/types/condition"
	"github.com/plarun/scheduler/api/types/entity/task"
)

type stat map[string]task.State

// func initAllocator() {
// 	// export configs
// 	if err := config.LoadConfig(); err != nil {
// 		log.Fatal(err)
// 	}

// 	// connect to mysql db
// 	mysql.ConnectDB()
// }

func TestConditionCheck(t *testing.T) {
	tests := map[string]struct {
		condition string
		stat      stat
		want      bool
	}{
		"good 1": {condition: "su(task1)", stat: stat{"task1": task.StateSuccess}, want: true},
		"good 2": {condition: "su(task2)", stat: stat{"task2": task.StateSuccess}, want: true},
		"good 3": {condition: "su(task3)", stat: stat{"task3": task.StateSuccess}, want: true},
		"good 4": {
			condition: "su(task1)&su(task2)&(fa(task3)|nr(task4))",
			stat: stat{
				"task1": task.StateSuccess, "task2": task.StateSuccess,
				"task3": task.StateFailure, "task4": task.StateFailure,
			}, want: true},
		"good 5": {
			condition: "su(task1)&su(task2)&su(task3)&su(task4)",
			stat: stat{
				"task1": task.StateSuccess, "task2": task.StateSuccess,
				"task3": task.StateSuccess, "task4": task.StateSuccess,
			}, want: true},
		"good 6": {
			condition: "su(task1)&su(task2)&su(task3)&su(task4)",
			stat: stat{
				"task1": task.StateSuccess, "task2": task.StateSuccess,
				"task3": task.StateSuccess, "task4": task.StateFailure,
			}, want: false},
		"good 7": {
			condition: "su(task1)&su(task2)&su(task3)&su(task4)",
			stat: stat{
				"task1": task.StateFailure, "task2": task.StateSuccess,
				"task3": task.StateSuccess, "task4": task.StateSuccess,
			}, want: false},
		"good 8": {condition: "", stat: stat{}, want: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			chk := NewConditionChecker(0)
			chk.initiated = true
			chk.cond = tc.condition
			chk.condTaskStatus = tc.stat

			if expr, err := condition.Build(chk.cond); err != nil {
				t.Fatalf("Build: %v", err)
			} else {
				chk.expr = expr
			}

			if got, err := chk.Check(); err != nil {
				t.Fatalf("Check: %v", err)
			} else if got != tc.want {
				t.Fatalf("want: %#v, got: %#v", tc.want, got)
			}
		})
	}
}
