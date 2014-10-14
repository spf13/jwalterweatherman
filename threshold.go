package jww

type Threshold int

const (
	ThresholdTrace Threshold = iota
	ThresholdDebug
	ThresholdInfo
	ThresholdWarn
	ThresholdError
	ThresholdCritical
	ThresholdFatal
	MaxThreshold
)

var thresholdPrefixes map[Threshold]string = map[Threshold]string{
	ThresholdTrace: "TRACE ",
	ThresholdDebug: "DEBUG ",
	ThresholdInfo: "INFO ",
	ThresholdWarn: "WARN ",
	ThresholdError: "ERROR ",
	ThresholdCritical: "CRITICAL ",
	ThresholdFatal: "FATAL ",
}
