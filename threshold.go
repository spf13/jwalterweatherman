package jww

type Threshold int

const (
	Trace Threshold = iota
	Debug
	Info
	Warn
	Error
	Critical
	Fatal
)

var prefixes map[Threshold]string = map[Threshold]string{
	Trace:    "TRACE ",
	Debug:    "DEBUG ",
	Info:     "INFO ",
	Warn:     "WARN ",
	Error:    "ERROR ",
	Critical: "CRITICAL ",
	Fatal:    "FATAL ",
}
