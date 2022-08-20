package builder

import (
	"fmt"
	"regexp"
)

// checkTimeFormat checks the format of time as hh:mm:ss
func parseTime(time string) error {
	// regex to check time string in format hh:mm:ss
	timeRegex, _ := regexp.Compile("^([0-1]?[0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$")
	if ok := timeRegex.MatchString(time); ok {
		return nil
	}
	return fmt.Errorf("start_time: %v is not allowed, only hh:mm:ss format is allowed", time)
}
