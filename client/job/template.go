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

func runStatus(jobName string, startTime string, endTime string, status string) {
	fmt.Println(line1)
	fmt.Printf("| %-62s | %-19s | %-19s | %-7s |\n", "Job Name", "Start Time", "End Time", "Status")
	fmt.Println(line1)
	fmt.Printf("| %-62s | %-19s | %-19s | %-7s |\n", jobName, startTime, endTime, status)
	fmt.Println(line1)
}

// func runHistory(jobName string, startTimes []string, endTimes []string, status []string) {
// 	for i := 0; i < len(startTimes); i++ {
// 		runStatus(jobName, startTimes[i], endTimes[i], status[i])
// 	}
// }

func jobDefinition(jobDef *pb.GetJilRes) {
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
