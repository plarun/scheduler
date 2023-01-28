package service

import (
	"reflect"
	"testing"

	"github.com/plarun/scheduler/api/types/condition"
	"github.com/plarun/scheduler/api/types/entity/task"
	"github.com/plarun/scheduler/internal/validator/errors"
	"github.com/plarun/scheduler/proto"
)

func TestTaskName(t *testing.T) {
	tests := map[string]struct {
		input   string
		wantErr error
	}{
		"good job name":    {input: "test_job_1"},
		"bad space":        {input: "test_ job", wantErr: errors.ErrTaskInvalidChar},
		"bad special char": {input: "test_#", wantErr: errors.ErrTaskInvalidChar},
		"bad long name":    {input: "test_job_12345678901234567890123456789012345678901234567890123456", wantErr: errors.ErrTaskMaxLength},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := checkFieldTaskName(tc.input)
			if tc.wantErr != nil {
				if tc.wantErr != err {
					t.Fatalf("want: %#v, got: %#v", tc.wantErr, err)
				}
			} else if err != nil {
				t.Fatal("failed")
			}
		})
	}
}

func TestCastFieldAction(t *testing.T) {
	tests := map[string]struct {
		input   string
		want    task.Action
		wantErr error
	}{
		"insert action":        {input: "insert_task", want: task.ActionInsert},
		"update action":        {input: "update_task", want: task.ActionUpdate},
		"delete action":        {input: "delete_task", want: task.ActionDelete},
		"bad empty action":     {input: "", wantErr: errors.ErrNonEmptyValueRequired},
		"bad case sensitive 1": {input: "Insert_task", wantErr: errors.ErrInvalidActionAttr},
		"bad case sensitive 2": {input: "INSERT_task", wantErr: errors.ErrInvalidActionAttr},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			g, err := castFieldAction(tc.input)
			got := task.Action(g)
			if tc.wantErr != nil {
				if tc.wantErr != err {
					t.Fatalf("want: %#v, got: %#v", tc.wantErr, err)
				}
			} else {
				if !reflect.DeepEqual(tc.want, got) {
					t.Fatalf("want: %#v, got: %#v", tc.want, got)
				}
			}
		})
	}
}

func TestCastFieldType(t *testing.T) {
	tests := map[string]struct {
		input   string
		want    task.Type
		wantErr error
	}{
		"callable type":        {input: "callable", want: task.TypeCallable, wantErr: nil},
		"bundle type":          {input: "bundle", want: task.TypeBundle, wantErr: nil},
		"bad empty":            {input: "", wantErr: errors.ErrNonEmptyValueRequired},
		"bad case sensitive 1": {input: "Callable", wantErr: errors.ErrInvalidTypeAttr},
		"bad case sensitive 2": {input: "CALLABLE", wantErr: errors.ErrInvalidTypeAttr},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			v, err := castFieldType(tc.input)
			got := task.Type(v)
			if tc.wantErr != nil {
				if tc.wantErr != err {
					t.Fatalf("want: %#v, got: %#v", tc.wantErr, err)
				}
			} else {
				if !reflect.DeepEqual(tc.want, got) {
					t.Fatalf("want: %#v, got: %#v", tc.want, got)
				}
			}
		})
	}
}

func TestCastFieldStartCondition(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
		fail  bool
	}{
		// negative cases
		"bad 1": {input: "su(job1", fail: true},
		"bad 2": {input: "su(job1 & su(job2)", fail: true},
		"bad 3": {input: "su(job1) | su()", fail: true},
		"bad 4": {input: "s(job1)", fail: true},
		"bad 5": {input: "(su(job1) & su(job2)", fail: true},
		"bad 6": {input: "su(job1) & su(job2))", fail: true},
		"bad 7": {input: "suc(job1)", fail: true},
		"bad 8": {input: "su(job1) & suc(job2)", fail: true},
		"bad 9": {input: "su(job1) &", fail: true},

		// positive cases
		"good 1":  {input: "", want: ""},
		"good 2":  {input: "su(job1)", want: "su(job1)"},
		"good 3":  {input: "fa(job1)", want: "fa(job1)"},
		"good 4":  {input: "nr(job1)", want: "nr(job1)"},
		"good 5":  {input: "", want: ""},
		"good 6":  {input: "(su(job1))", want: "(su(job1))"},
		"good 7":  {input: "su(job1) & su(job2)", want: "su(job1)&su(job2)"},
		"good 8":  {input: "su(job1) & fa(job2)", want: "su(job1)&fa(job2)"},
		"good 9":  {input: "(su(job1) & su(job2))", want: "(su(job1)&su(job2))"},
		"good 10": {input: "su(job1) & fa(job2) | nr(job3)", want: "su(job1)&fa(job2)|nr(job3)"},
		"good 11": {input: "(su(job1) & su(job2)) | su(job3)", want: "(su(job1)&su(job2))|su(job3)"},
		"good 12": {input: "((su(job1) & su(job2)) | su(job3))", want: "((su(job1)&su(job2))|su(job3))"},
		"good 13": {input: "su(job1) & (su(job2) | su(job3))", want: "su(job1)&(su(job2)|su(job3))"},
		"good 14": {input: "(su(job1) & (su(job2) | su(job3)))", want: "(su(job1)&(su(job2)|su(job3)))"},
		"good 15": {input: "su(job1) & (su(job2) | su(job3) & (fa(job4) | nr(job5)))", want: "su(job1)&(su(job2)|su(job3)&(fa(job4)|nr(job5)))"},
		"good 16": {input: "su(job1)|(su(job2)&su(job3))|(fa(job4)|nr(job5))", want: "su(job1)|(su(job2)&su(job3))|(fa(job4)|nr(job5))"},
		"good 17": {input: "su(job1) | (((su(job2) & su(job3) | fa(job4)) | nr(job5)) | nr(job6)) & su(job7)", want: "su(job1)|(((su(job2)&su(job3)|fa(job4))|nr(job5))|nr(job6))&su(job7)"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			clause, err := condition.Build(tc.input)
			if tc.fail {
				if err == nil {
					t.Fatalf("should fail")
				}
			} else {
				got := clause.String()
				if !reflect.DeepEqual(tc.want, got) {
					t.Fatalf("want: %#v, got: %#v", tc.want, got)
				}
			}
		})
	}
}

func TestCastFieldRundays(t *testing.T) {
	var (
		SU task.RunDay = task.RunDaySunday
		MO task.RunDay = task.RunDayMonday
		TU task.RunDay = task.RunDayTuesday
		WE task.RunDay = task.RunDayWednesday
		TH task.RunDay = task.RunDayThursday
		FR task.RunDay = task.RunDayFriday
		SA task.RunDay = task.RunDaySaturday

		NO_DAY   task.RunDay = task.RunDayEmpty
		ALL_DAYS task.RunDay = task.RunDayDaily
	)
	tests := map[string]struct {
		input   string
		want    task.RunDay
		wantErr error
	}{
		"good 1":  {input: "", want: NO_DAY},
		"good 2":  {input: "su", want: SU},
		"good 3":  {input: "mo", want: MO},
		"good 4":  {input: "tu", want: TU},
		"good 5":  {input: "we", want: WE},
		"good 6":  {input: "th", want: TH},
		"good 7":  {input: "fr", want: FR},
		"good 8":  {input: "sa", want: SA},
		"good 9":  {input: "mo,tu", want: MO | TU},
		"good 10": {input: "mo,sa,su", want: SU | MO | SA},
		"good 11": {input: "mo,tu,we,th,fr,sa,su", want: ALL_DAYS},
		"good 12": {input: "sa,fr,th,we,tu,mo,su", want: ALL_DAYS},
		"good 13": {input: "sa,fr,th,we,tu,mo,su", want: SU | MO | TU | WE | TH | FR | SA},
		"good 14": {input: "mo, tu, we", want: MO | TU | WE},
		"good 15": {input: "su,mo, tu,we, th", want: SU | MO | TU | WE | TH},
		"good 16": {input: "SU", want: SU},
		"good 17": {input: "su, Mo", want: SU | MO},
		"good 18": {input: "all", want: ALL_DAYS},
		"good 19": {input: "All", want: ALL_DAYS},
		"good 20": {input: "All", want: ALL_DAYS},

		"bad 1": {input: "su,mo,su", wantErr: errors.ErrRepeatedRundaysAttr},
		"bad 2": {input: "su,mo,", wantErr: errors.ErrInvalidRundaysAttr},
		"bad 3": {input: ",su,mo,tu", wantErr: errors.ErrInvalidRundaysAttr},
		"bad 4": {input: ",su,,tu", wantErr: errors.ErrInvalidRundaysAttr},
		"bad 5": {input: "su,mo,all", wantErr: errors.ErrInvalidRundaysAttr},
		"bad 6": {input: "all,", wantErr: errors.ErrInvalidRundaysAttr},
		"bad 7": {input: ",all,", wantErr: errors.ErrInvalidRundaysAttr},
		"bad 8": {input: ",", wantErr: errors.ErrInvalidRundaysAttr},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			g, err := castFieldRundays(tc.input)
			got := task.RunDay(g)
			if tc.wantErr != nil {
				if tc.wantErr != err {
					t.Fatalf("want: %#v, got: %#v", tc.wantErr, err)
				}
			} else {
				if !reflect.DeepEqual(tc.want, got) {
					t.Fatalf("want: %#v, got: %#v", tc.want, got)
				}
			}
		})
	}
}

func TestCastStartTimes(t *testing.T) {
	tests := map[string]struct {
		input   string
		want    []string
		wantErr error
	}{
		"good 1": {input: "00:00", want: []string{"00:00"}},
		"good 2": {input: "23:59", want: []string{"23:59"}},
		"good 3": {input: "0:00,1:00,02:00", want: []string{"00:00", "01:00", "02:00"}},
		"good 4": {input: "0:00,1:00, 02:00", want: []string{"00:00", "01:00", "02:00"}},
		"good 5": {input: "12:00, 13:30:30", want: []string{"12:00", "13:30:30"}},

		"bad 1": {input: "", wantErr: errors.ErrNonEmptyValueRequired},
		"bad 2": {input: "0:00,1:00,", wantErr: errors.ErrInvalidStartTimesAttr},
		"bad 3": {input: "0:00, 24:00", wantErr: errors.ErrInvalidStartTimesAttr},
		"bad 4": {input: "12:30:60", wantErr: errors.ErrInvalidStartTimesAttr},
		"bad 5": {input: "12:5", wantErr: errors.ErrInvalidStartTimesAttr},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := castFieldStartTimes(tc.input)
			if tc.wantErr != nil {
				if tc.wantErr != err {
					t.Fatalf("want: %#v, got: %#v", tc.wantErr, err)
				}
			} else {
				if !reflect.DeepEqual(tc.want, got) {
					t.Fatalf("want: %#v, got: %#v", tc.want, got)
				}
			}
		})
	}
}

func TestCastStartMins(t *testing.T) {
	tests := map[string]struct {
		input   string
		want    []uint8
		wantErr error
	}{
		"good 1": {input: "0", want: []uint8{0}},
		"good 2": {input: "00", want: []uint8{0}},
		"good 3": {input: "10,20,30", want: []uint8{10, 20, 30}},

		"bad 1": {input: "", wantErr: errors.ErrNonEmptyValueRequired},
		"bad 2": {input: "60", wantErr: errors.ErrInvalidStartMinsAttr},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := castFieldStartMins(tc.input)
			if tc.wantErr != nil {
				if tc.wantErr != err {
					t.Fatalf("want: %#v, got: %#v", tc.wantErr, err)
				}
			} else {
				if !reflect.DeepEqual(tc.want, got) {
					t.Fatalf("want: %#v, got: %#v", tc.want, got)
				}
			}
		})
	}
}

func TestCastRunwindow(t *testing.T) {
	tests := map[string]struct {
		input   string
		want    *proto.TimeRange
		wantErr error
	}{
		"good 1": {input: "00:00-2:00", want: &proto.TimeRange{Start: "00:00", End: "02:00"}},
		"good 2": {input: "10:30-22:30", want: &proto.TimeRange{Start: "10:30", End: "22:30"}},
		"good 3": {input: "00:10-02:50", want: &proto.TimeRange{Start: "00:10", End: "02:50"}},
		"good 4": {input: "4:00-8:00", want: &proto.TimeRange{Start: "04:00", End: "08:00"}},

		"bad 1": {input: "", wantErr: errors.ErrNonEmptyValueRequired},
		"bad 2": {input: "12:00", wantErr: errors.ErrInvalidRunWindowAttr},
		"bad 3": {input: "0:00,1:00", wantErr: errors.ErrInvalidRunWindowAttr},
		"bad 4": {input: "02:00-03:00-04:00", wantErr: errors.ErrInvalidRunWindowAttr},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			start, end, err := castFieldRunWindow(tc.input)
			if tc.wantErr != nil {
				if tc.wantErr != err {
					t.Fatalf("want: %#v, got: %#v", tc.wantErr, err)
				}
			} else {
				if !reflect.DeepEqual(tc.want.Start, start) || !reflect.DeepEqual(tc.want.End, end) {
					t.Fatalf("want: %#v, got: %#v - %#v", tc.want, start, end)
				}
			}
		})
	}
}

func TestCastPriority(t *testing.T) {
	tests := map[string]struct {
		input   string
		want    int32
		wantErr error
	}{
		"good 1": {input: "0", want: 0},
		"good 2": {input: "1", want: 1},
		"good 3": {input: "2", want: 2},
		"good 4": {input: "3", want: 3},
		"good 5": {input: "low", want: 0},
		"good 6": {input: "normal", want: 1},
		"good 7": {input: "important", want: 2},
		"good 8": {input: "critical", want: 3},
		"good 9": {input: "", want: 0},

		"bad 1": {input: "5", wantErr: errors.ErrInvalidPriorityAttr},
		"bad 2": {input: "Low", wantErr: errors.ErrInvalidPriorityAttr},
		"bad 3": {input: "LOW", wantErr: errors.ErrInvalidPriorityAttr},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := castFieldPriority(tc.input)
			if tc.wantErr != nil {
				if tc.wantErr != err {
					t.Fatalf("want: %#v, got: %#v", tc.wantErr, err)
				}
			} else {
				if !reflect.DeepEqual(tc.want, got) {
					t.Fatalf("want: %#v, got: %#v", tc.want, got)
				}
			}
		})
	}
}
