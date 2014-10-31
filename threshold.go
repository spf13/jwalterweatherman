package jww

type Threshold int

const (
	LevelTrace Threshold = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelCritical
	LevelFatal
)

var prefixes map[Threshold]string = map[Threshold]string{
	LevelTrace:    "TRACE ",
	LevelDebug:    "DEBUG ",
	LevelInfo:     "INFO ",
	LevelWarn:     "WARN ",
	LevelError:    "ERROR ",
	LevelCritical: "CRITICAL ",
	LevelFatal:    "FATAL ",
}
