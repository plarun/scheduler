package job

import (
	"fmt"
	"strings"

	pb "github.com/plarun/scheduler/client/data"
)

const (
	line1 = "+----------------------------------------------------------------+---------------------+---------------------+---------+"
	line2 = "+----------------+----------------------------------------------------------------------------------------------+"
)

// printRunStatus prints the last run status of job
func printRunStatus(jobName string, startTime string, endTime string, status string) {
	fmt.Println(line1)
	fmt.Printf("| %-62s | %-19s | %-19s | %-7s |\n", "Job Name", "Start Time", "End Time", "Status")
	fmt.Println(line1)
	fmt.Printf("| %-62s | %-19s | %-19s | %-7s |\n", jobName, startTime, endTime, status)
	fmt.Println(line1)
}

// printRunHistory prints the previous run history of job
func printRunHistory(jobName string, startTimes []string, endTimes []string, status []pb.Status) {
	fmt.Println(line1)
	fmt.Printf("| %-62s | %-19s | %-19s | %-7s |\n", "Job Name", "Start Time", "End Time", "Status")
	fmt.Println(line1)
	for i := 0; i < len(startTimes); i++ {
		fmt.Printf("| %-62s | %-19s | %-19s | %-7s |\n", jobName, startTimes[i], endTimes[i], status[i].String())
	}
	fmt.Println(line1)
}

// printJobDefinition prints the job definition
func printJobDefinition(jobDef *pb.GetJilRes) {
	fmt.Println(line2)
	fmt.Printf("| %-14s | %-92s |\n", "Job Name", jobDef.GetJobName())
	fmt.Printf("| %-14s | %-92s |\n", "Command", jobDef.GetCommand())
	fmt.Printf("| %-14s | %-92s |\n", "Conditions", strings.Join(jobDef.GetConditions(), " & "))
	fmt.Printf("| %-14s | %-92s |\n", "Std Out Log", jobDef.GetStdOut())
	fmt.Printf("| %-14s | %-92s |\n", "Std Err Log", jobDef.GetStdErr())
	fmt.Printf("| %-14s | %-92s |\n", "Machine", jobDef.GetMachine())
	fmt.Printf("| %-14s | %-92s |\n", "Run Days", jobDef.GetRunDays())
	fmt.Printf("| %-14s | %-92s |\n", "Start Times", jobDef.GetStartTimes())
	fmt.Println(line2)
}
