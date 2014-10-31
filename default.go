package jww

import (
	"io"
	"io/ioutil"
	"log"
)

var (
	TRACE    *log.Logger
	DEBUG    *log.Logger
	INFO     *log.Logger
	WARN     *log.Logger
	ERROR    *log.Logger
	CRITICAL *log.Logger
	FATAL    *log.Logger

	defaultNotepad *Notepad
)

func reloadDefaultNotepad() {
	TRACE = defaultNotepad.TRACE
	DEBUG = defaultNotepad.DEBUG
	INFO = defaultNotepad.INFO
	WARN = defaultNotepad.WARN
	ERROR = defaultNotepad.ERROR
	CRITICAL = defaultNotepad.CRITICAL
	FATAL = defaultNotepad.FATAL
}

func init() {
	defaultNotepad = NewNotepad(LevelInfo, LevelTrace, ioutil.Discard, "", log.Ldate|log.Ltime)
	reloadDefaultNotepad()
}

// SetLogThreshold set the log threshold for the default notepad. Trace by default.
func SetLogThreshold(threshold Threshold) {
	defaultNotepad.SetLogThreshold(threshold)
	reloadDefaultNotepad()
}

// SetLogOutput set the log output for the default notepad. Discarded by default.
func SetLogOutput(handle io.Writer) {
	defaultNotepad.SetLogOutput(handle)
	reloadDefaultNotepad()
}

// SetStdoutThreshold set the standard output threshold for the default notepad.
// Info by default.
func SetStdoutThreshold(threshold Threshold) {
	defaultNotepad.SetStdoutThreshold(threshold)
	reloadDefaultNotepad()
}

// SetPrefix set the prefix for the default notepad. Empty by default.
func SetPrefix(prefix string) {
	defaultNotepad.SetPrefix(prefix)
	reloadDefaultNotepad()
}

// SetFlags set the flags for the default notepad. "log.Ldate | log.Ltime" by default.
func SetFlags(flags int) {
	defaultNotepad.SetFlags(flags)
	reloadDefaultNotepad()
}
