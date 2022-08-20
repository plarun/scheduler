package builder

import (
	"fmt"
	"github.com/plarun/scheduler/client/model"
	"strings"
)

// validRunDays checks format of run days
func validRunDays(days []string) error {
	var daysBitFlag int8
	// remove white spaces or tabs between week day
	// convert the week days to lower case
	for i := 0; i < len(days); i++ {
		days[i] = strings.ToLower(strings.Trim(days[i], " \t"))
	}
	for _, day := range days {
		repeatedDay := false
		switch day {
		case "su":
			if daysBitFlag&model.SU != 0 {
				repeatedDay = true
			} else {
				daysBitFlag = model.SU | daysBitFlag
			}
		case "mo":
			if daysBitFlag&model.MO != 0 {
				repeatedDay = true
			} else {
				daysBitFlag = model.MO | daysBitFlag
			}
		case "tu":
			if daysBitFlag&model.TU != 0 {
				repeatedDay = true
			} else {
				daysBitFlag = model.TU | daysBitFlag
			}
		case "we":
			if daysBitFlag&model.WE != 0 {
				repeatedDay = true
			} else {
				daysBitFlag = model.WE | daysBitFlag
			}
		case "th":
			if daysBitFlag&model.TH != 0 {
				repeatedDay = true
			} else {
				daysBitFlag = model.TH | daysBitFlag
			}
		case "fr":
			if daysBitFlag&model.FR != 0 {
				repeatedDay = true
			} else {
				daysBitFlag = model.FR | daysBitFlag
			}
		case "sa":
			if daysBitFlag&model.SA != 0 {
				repeatedDay = true
			} else {
				daysBitFlag = model.SA | daysBitFlag
			}
		default:
			return fmt.Errorf("invalid week: %v, days should be one of su,mo,tu,we,th,fr,sa", day)
		}
		if repeatedDay {
			return fmt.Errorf("week: %v is repated", day)
		}
	}
	return nil
}
