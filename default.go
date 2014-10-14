package jww

import (
	"io"
	"io/ioutil"
	"log"
)

var (
	TRACE      *log.Logger
	DEBUG      *log.Logger
	INFO       *log.Logger
	WARN       *log.Logger
	ERROR      *log.Logger
	CRITICAL   *log.Logger
	FATAL      *log.Logger
	
	defaultNotepad *Notepad
)

func init () {
	defaultNotepad = NewNotepad(ThresholdInfo, ThresholdTrace, ioutil.Discard, "", log.Ldate | log.Ltime)
	reloadDefaultNotepad()
}

func reloadDefaultNotepad () {
	TRACE = defaultNotepad.TRACE
	DEBUG = defaultNotepad.DEBUG
	INFO = defaultNotepad.INFO
	WARN = defaultNotepad.WARN
	ERROR = defaultNotepad.ERROR
	CRITICAL = defaultNotepad.CRITICAL
	FATAL = defaultNotepad.FATAL
}

func SetLogThreshold (threshold Threshold) {
	defaultNotepad.SetLogThreshold(threshold)
	reloadDefaultNotepad()
}

func SetLogOutput (handle io.Writer) {
	defaultNotepad.SetLogOutput(handle)
	reloadDefaultNotepad()
}

func SetStdoutThreshold (threshold Threshold) {
	defaultNotepad.SetStdoutThreshold(threshold)
	reloadDefaultNotepad()
}

func SetPrefix (prefix string) {
	defaultNotepad.SetPrefix(prefix)
	reloadDefaultNotepad()
}

func SetFlags (flags int) {
	defaultNotepad.SetFlags(flags)
	reloadDefaultNotepad()
}
