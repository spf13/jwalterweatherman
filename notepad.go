package jww

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

// Notepad is where you leave a note !
type Notepad struct {
	TRACE      *log.Logger
	DEBUG      *log.Logger
	INFO       *log.Logger
	WARN       *log.Logger
	ERROR      *log.Logger
	CRITICAL   *log.Logger
	FATAL      *log.Logger
	
	loggers    []**log.Logger
	logHandle  io.Writer
	logThreshold   Threshold
	stdoutThreshold   Threshold
	prefix     string
	flags      int
}

// NewNotepad create a new notepad.
func NewNotepad(stdoutThreshold Threshold, logThreshold Threshold, logHandle io.Writer, prefix string, flags int) *Notepad {
	n := &Notepad{}
	
	n.loggers = append(n.loggers, &n.TRACE, &n.DEBUG, &n.INFO, &n.WARN, &n.ERROR, &n.CRITICAL, &n.FATAL)
	n.logHandle = logHandle
	n.logThreshold = logThreshold
	n.stdoutThreshold = stdoutThreshold
	
	if len(prefix) != 0 {
		n.prefix = "["+prefix+"] "
	} else {
		n.prefix = ""
	}
	
	n.flags = flags
	
	n.init()
	
	return n
}

// init create the loggers for each level depending on the notepad thresholds
func (n *Notepad) init () {
	bothHandle := io.MultiWriter(os.Stdout, n.logHandle)
	
	for t, logger := range n.loggers {
		threshold := Threshold(t)
		switch {
			case threshold >= n.logThreshold && threshold >= n.stdoutThreshold:
				*logger = log.New(bothHandle, n.prefix + prefixes[threshold], n.flags)
				
			case threshold >= n.logThreshold:
				*logger = log.New(n.logHandle, n.prefix + prefixes[threshold], n.flags)
				
			case threshold >= n.stdoutThreshold:
				*logger = log.New(os.Stdout, n.prefix + prefixes[threshold], n.flags)
				
			default:
				*logger = log.New(ioutil.Discard, n.prefix + prefixes[threshold], n.flags)
		}
	}
}

// SetLogThreshold change the threshold above which messages are written to the
// log file
func (n *Notepad) SetLogThreshold (threshold Threshold) {
	n.logThreshold = threshold
	n.init()
}

// SetLogOutput change the file where log messages are written
func (n *Notepad) SetLogOutput (handle io.Writer) {
	n.logHandle = handle
	n.init()
}

// SetStdoutThreshold change the threshold above which messages are written to the
// standard output
func (n *Notepad) SetStdoutThreshold (threshold Threshold) {
	n.stdoutThreshold = threshold
	n.init()
}

// SetPrefix change the prefix used by the notepad. Prefixes are displayed between 
// brackets at the begining of the line. An empty prefix won't be displayed at all.
func (n *Notepad) SetPrefix (prefix string) {
	if len(prefix) != 0 {
		n.prefix = "["+prefix+"] "
	} else {
		n.prefix = ""
	}
	n.init()
}

// SetFlags choose which flags the logger will display (after prefix and message
// level). See the package log for more informations on this.
func (n *Notepad) SetFlags (flags int) {
	n.flags = flags
	n.init()
}
