package service

import (
	"log"
	"reflect"
	"testing"

	"github.com/plarun/scheduler/config"
	"github.com/plarun/scheduler/internal/validator/db/mysql"
	"github.com/plarun/scheduler/proto"
)

func initValidator() {
	// export configs
	if err := config.LoadConfig(); err != nil {
		log.Fatal(err)
	}

	// connect to mysql db
	mysql.ConnectDB()
}

func compareTaskEntity(t1, t2 *proto.ValidatedTaskEntity) bool {
	if !reflect.DeepEqual(t1.Action, t2.Action) ||
		!reflect.DeepEqual(t1.Name, t2.Name) ||
		!reflect.DeepEqual(t1.Type, t2.Type) ||
		!reflect.DeepEqual(t1.Parent, t2.Parent) ||
		!reflect.DeepEqual(t1.Machine, t2.Machine) ||
		!reflect.DeepEqual(t1.Command, t2.Command) ||
		!reflect.DeepEqual(t1.Condition, t2.Condition) ||
		!reflect.DeepEqual(t1.OutLogFile, t2.OutLogFile) ||
		!reflect.DeepEqual(t1.ErrLogFile, t2.ErrLogFile) ||
		!reflect.DeepEqual(t1.Label, t2.Label) ||
		!reflect.DeepEqual(t1.Profile, t2.Profile) ||
		!reflect.DeepEqual(t1.RunDays, t2.RunDays) ||
		!reflect.DeepEqual(t1.StartTimes, t2.StartTimes) ||
		!reflect.DeepEqual(t1.RunWindow, t2.RunWindow) ||
		!reflect.DeepEqual(t1.StartMins, t2.StartMins) ||
		!reflect.DeepEqual(t1.Priority, t2.Priority) {
		return true
	}
	return false
}

func TestInsertJob(t *testing.T) {
	initValidator()

	tests := map[string]struct {
		input   *proto.ParsedTaskEntity
		want    *proto.ValidatedTaskEntity
		wantErr error
	}{
		"good job name": {
			input: &proto.ParsedTaskEntity{
				Action: "insert_task",
				Target: "test_1_job_1",
				Fields: map[string]string{
					"type":         "callable",
					"run_days":     "mo,we,th,fr",
					"start_times":  "11:00, 12:00, 13:00",
					"label":        "callable batch job",
					"command":      "/opt/work/scripts/test.sh 5",
					"out_log_file": "/opt/work/logs/test_1_job_1.out",
					"err_log_file": "/opt/work/logs/test_1_job_1.err",
				},
			},
			want: &proto.ValidatedTaskEntity{
				Action:     "insert_task",
				Name:       "test_1_job_1",
				Type:       &proto.NullableString{Value: "callable", Flag: proto.NullableFlag_Available},
				Parent:     &proto.NullableString{Value: "", Flag: proto.NullableFlag_NotAvailable},
				Machine:    &proto.NullableString{Value: "", Flag: proto.NullableFlag_NotAvailable},
				Command:    &proto.NullableString{Value: "/opt/work/scripts/test.sh 5", Flag: proto.NullableFlag_Available},
				Condition:  &proto.NullableString{Value: "", Flag: proto.NullableFlag_NotAvailable},
				OutLogFile: &proto.NullableString{Value: "/opt/work/logs/test_1_job_1.out", Flag: proto.NullableFlag_Available},
				ErrLogFile: &proto.NullableString{Value: "/opt/work/logs/test_1_job_1.err", Flag: proto.NullableFlag_Available},
				Label:      &proto.NullableString{Value: "callable batch job", Flag: proto.NullableFlag_Available},
				Profile:    &proto.NullableString{Value: "", Flag: proto.NullableFlag_NotAvailable},
				RunDays:    &proto.NullableInt32{Value: 0, Flag: proto.NullableFlag_NotAvailable},
				StartTimes: &proto.NullableStrings{Value: []string{"11:00", "12:00", "13:00"}, Flag: proto.NullableFlag_Available},
				RunWindow:  &proto.NullableTimeRange{Value: nil, Flag: proto.NullableFlag_NotAvailable},
				StartMins:  &proto.NullableInt32S{Value: []int32{}, Flag: proto.NullableFlag_NotAvailable},
				Priority:   &proto.NullableInt32{Value: 0, Flag: proto.NullableFlag_NotAvailable},
			},
		},
		// "bad space":        {input: "test_ job", wantErr: errors.ErrJobInvalidChar},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ValidateTaskAction(tc.input)
			if tc.wantErr != nil {
				if tc.wantErr != err {
					t.Fatalf("want: %#v, got: %#v", tc.wantErr, err)
				}
			} else {
				// if !reflect.DeepEqual(tc.want, got) {
				if !compareTaskEntity(tc.want, got) {
					t.Fatalf("want: %#v, got: %#v", tc.want, got)
				}
			}
		})
	}
}
