package jww

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

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

func (n *Notepad) init () {
	bothHandle := io.MultiWriter(os.Stdout, n.logHandle)
	
	for t, logger := range n.loggers {
		threshold := Threshold(t)
		switch {
			case threshold >= n.logThreshold && threshold >= n.stdoutThreshold:
				*logger = log.New(bothHandle, n.prefix + thresholdPrefixes[threshold], n.flags)
				
			case threshold >= n.logThreshold:
				*logger = log.New(n.logHandle, n.prefix + thresholdPrefixes[threshold], n.flags)
				
			case threshold >= n.stdoutThreshold:
				*logger = log.New(os.Stdout, n.prefix + thresholdPrefixes[threshold], n.flags)
				
			default:
				*logger = log.New(ioutil.Discard, n.prefix + thresholdPrefixes[threshold], n.flags)
		}
	}
}

func (n *Notepad) SetLogThreshold (threshold Threshold) {
	n.logThreshold = threshold
	n.init()
}

func (n *Notepad) SetLogOutput (handle io.Writer) {
	n.logHandle = handle
	n.init()
}

func (n *Notepad) SetStdoutThreshold (threshold Threshold) {
	n.stdoutThreshold = threshold
	n.init()
}

func (n *Notepad) SetPrefix (prefix string) {
	if len(prefix) != 0 {
		n.prefix = "["+prefix+"] "
	} else {
		n.prefix = ""
	}
	n.init()
}

func (n *Notepad) SetFlags (flags int) {
	n.flags = flags
	n.init()
}
