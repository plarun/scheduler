package time

const (
	fmtHHMM   string = "15:04"
	fmtHHMMSS string = "15:04:05"

	fmtyymmdd         string = "060102"
	fmtyyddmm         string = "060201"
	fmtyyyymmdd       string = "20060102"
	fmtyyyyddmm       string = "20060201"
	fmtyyyymmddHHMM   string = fmtyyyymmdd + fmtHHMM
	fmtyyyymmddHHMMSS string = fmtyyyymmdd + fmtHHMMSS
	fmtyyyyddmmHHMM   string = fmtyyyyddmm + fmtHHMM
	fmtyyyyddmmHHMMSS string = fmtyyyyddmm + fmtHHMMSS

	default_Time     string = fmtHHMM
	default_Date     string = fmtyyyymmdd
	default_DateTime string = default_Date + default_Time
)

var formats map[string]string = map[string]string{
	"HHMM":           fmtHHMM,
	"HHMMSS":         fmtHHMMSS,
	"yymmdd":         fmtyymmdd,
	"yyddmm":         fmtyyddmm,
	"yyyymmdd":       fmtyyyymmdd,
	"yyyyddmm":       fmtyyyyddmm,
	"yyyymmddHHMM":   fmtyyyymmddHHMM,
	"yyyymmddHHMMSS": fmtyyyymmddHHMMSS,
	"yyyyddmmHHMM":   fmtyyyyddmmHHMM,
	"yyyyddmmHHMMSS": fmtyyyyddmmHHMMSS,
}

func GetLayout(format string) (string, bool) {
	lay, ok := formats[format]
	return lay, ok
}

func GetDefaultDateLayout() string {
	return default_Date
}

func GetDefaultTimeLayout() string {
	return default_Time
}

func GetDefaultDateTimeLayout() string {
	return default_DateTime
}
