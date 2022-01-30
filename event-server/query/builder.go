package query

import (
	"fmt"
	"strings"

	pb "github.com/plarun/scheduler/event-server/data"
	"github.com/plarun/scheduler/event-server/model"
)

func buildJobUpdateQuery(jobData *pb.Jil) string {
	var data *pb.JilData = jobData.Data

	var columns []string

	if jobData.AttributeFlag&model.COMMAND != 0 {
		columns = append(columns, fmt.Sprintf("command = '%s'", data.Command))
	}
	if jobData.AttributeFlag&model.STD_OUT != 0 {
		columns = append(columns, fmt.Sprintf("std_out_log = '%s'", data.StdOut))
	}
	if jobData.AttributeFlag&model.STD_ERR != 0 {
		columns = append(columns, fmt.Sprintf("std_err_log = '%s'", data.StdErr))
	}
	if jobData.AttributeFlag&model.MACHINE != 0 {
		columns = append(columns, fmt.Sprintf("machine = '%s'", data.Machine))
	}
	if jobData.AttributeFlag&model.START_TIMES != 0 {
		columns = append(columns, fmt.Sprintf("start_times = '%s'", data.StartTimes))
	}
	if jobData.AttributeFlag&model.RUN_DAYS != 0 {
		columns = append(columns, fmt.Sprintf("run_days = '%s'", data.RunDays))
	}

	columnStr := strings.Join(columns, ",")
	return columnStr
}

func buildJobStatusUpdateQuery(jobName string, status pb.Status) string {
	var columns []string

	columns = append(columns, fmt.Sprintf("status = '%s'", pb.Status_name[int32(status.Number())]))
	if status == pb.Status_RUNNING {
		columns = append(columns, "last_start_time = current_timestamp")
	} else if status == pb.Status_FAILED || status == pb.Status_SUCCESS {
		columns = append(columns, "last_end_time = current_timestamp")
	}

	columnStr := strings.Join(columns, ",")
	return columnStr
}
