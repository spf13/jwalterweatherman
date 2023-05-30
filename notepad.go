// Copyright © 2016 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package jwalterweatherman

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"github.com/fatih/color"
)

type Threshold int

func (t Threshold) String() string {
	return prefixes[t]
}

func (t Threshold) ColoredString() string {
	return coloredPrefixes[t]
}

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
	LevelTrace:    "TRACE",
	LevelDebug:    "DEBUG",
	LevelInfo:     "INFO",
	LevelWarn:     "WARN",
	LevelError:    "ERROR",
	LevelCritical: "CRITICAL",
	LevelFatal:    "FATAL",
}

var (
	DebugColor    = color.New(color.FgHiBlack)
	InfoColor     = color.New(color.FgHiBlue)
	WarnColor     = color.New(color.FgHiYellow)
	ErrorColor    = color.New(color.FgRed)
	CriticalColor = color.New(color.FgRed, color.Bold)
	FatalColor    = color.New(color.FgHiRed, color.Italic, color.Bold)
)

var coloredPrefixes map[Threshold]string = map[Threshold]string{
	LevelTrace:    "[TRACE]",
	LevelDebug:    fmt.Sprintf("[%s]", DebugColor.Sprint("DEBUG")),
	LevelInfo:     fmt.Sprintf("[%s]", InfoColor.Sprint("INFO")),
	LevelWarn:     fmt.Sprintf("[%s]", WarnColor.Sprint("WARN")),
	LevelError:    fmt.Sprintf("[%s]", ErrorColor.Sprint("ERROR")),
	LevelCritical: fmt.Sprintf("[%s]", CriticalColor.Sprint("CRITICAL")),
	LevelFatal:    fmt.Sprintf("[%s]", FatalColor.Sprint("FATAL")),
}

// Notepad is where you leave a note!
type Notepad struct {
	TRACE    *log.Logger
	DEBUG    *log.Logger
	INFO     *log.Logger
	WARN     *log.Logger
	ERROR    *log.Logger
	CRITICAL *log.Logger
	FATAL    *log.Logger

	LOG      *log.Logger
	FEEDBACK *Feedback

	loggers         [7]**log.Logger
	logHandle       io.Writer
	outHandle       io.Writer
	logThreshold    Threshold
	stdoutThreshold Threshold
	prefix          string
	flags           int
	colored         bool

	logListeners []LogListener
}

// A LogListener can be supplied to a Notepad to listen on log writes for a given
// threshold. This can be used to capture log events in unit tests and similar.
// Note that this function will be invoked once for each log threshold. If
// the given threshold is not of interest to you, return nil.
// Note that these listeners will receive log events for a given threshold, even
// if the current configuration says not to log it. That way you can count ERRORs even
// if you don't print them to the console.
type LogListener func(t Threshold) io.Writer

// NewNotepad creates a new Notepad.
func NewNotepad(
	outThreshold Threshold,
	logThreshold Threshold,
	outHandle, logHandle io.Writer,
	prefix string, flags int, color bool,
	logListeners ...LogListener,
) *Notepad {

	n := &Notepad{logListeners: logListeners}

	n.loggers = [7]**log.Logger{&n.TRACE, &n.DEBUG, &n.INFO, &n.WARN, &n.ERROR, &n.CRITICAL, &n.FATAL}
	n.outHandle = outHandle
	n.logHandle = logHandle
	n.stdoutThreshold = outThreshold
	n.logThreshold = logThreshold
	n.colored = color

	if len(prefix) != 0 {
		n.prefix = "[" + prefix + "] "
	} else {
		n.prefix = ""
	}

	n.flags = flags

	n.LOG = log.New(n.logHandle,
		"LOG:   ",
		n.flags)
	n.FEEDBACK = &Feedback{out: log.New(outHandle, "", 0), log: n.LOG}

	n.init()
	return n
}

// init creates the loggers for each level depending on the notepad thresholds.
func (n *Notepad) init() {
	logAndOut := io.MultiWriter(n.outHandle, n.logHandle)

	for t, logger := range n.loggers {
		threshold := Threshold(t)
		prefix := ""
		if n.colored {
			prefix = n.prefix + threshold.ColoredString() + " "
		} else {
			prefix = n.prefix + threshold.String() + " "
		}

		switch {
		case threshold >= n.logThreshold && threshold >= n.stdoutThreshold:
			*logger = log.New(n.createLogWriters(threshold, logAndOut), prefix, n.flags)

		case threshold >= n.logThreshold:
			*logger = log.New(n.createLogWriters(threshold, n.logHandle), prefix, n.flags)

		case threshold >= n.stdoutThreshold:
			*logger = log.New(n.createLogWriters(threshold, n.outHandle), prefix, n.flags)

		default:
			*logger = log.New(n.createLogWriters(threshold, ioutil.Discard), prefix, n.flags)
		}
	}
}

func (n *Notepad) createLogWriters(t Threshold, handle io.Writer) io.Writer {
	if len(n.logListeners) == 0 {
		return handle
	}
	writers := []io.Writer{handle}
	for _, l := range n.logListeners {
		w := l(t)
		if w != nil {
			writers = append(writers, w)
		}
	}

	if len(writers) == 1 {
		return handle
	}

	return io.MultiWriter(writers...)
}

// SetLogThreshold changes the threshold above which messages are written to the
// log file.
func (n *Notepad) SetLogThreshold(threshold Threshold) {
	n.logThreshold = threshold
	n.init()
}

// SetLogOutput changes the file where log messages are written.
func (n *Notepad) SetLogOutput(handle io.Writer) {
	n.logHandle = handle
	n.init()
}

// GetLogThreshold returns the defined Treshold for the log logger.
func (n *Notepad) GetLogThreshold() Threshold {
	return n.logThreshold
}

// SetStdoutThreshold changes the threshold above which messages are written to the
// standard output.
func (n *Notepad) SetStdoutThreshold(threshold Threshold) {
	n.stdoutThreshold = threshold
	n.init()
}

// GetStdoutThreshold returns the Treshold for the stdout logger.
func (n *Notepad) GetStdoutThreshold() Threshold {
	return n.stdoutThreshold
}

// SetPrefix changes the prefix used by the notepad. Prefixes are displayed between
// brackets at the beginning of the line. An empty prefix won't be displayed at all.
func (n *Notepad) SetPrefix(prefix string) {
	if len(prefix) != 0 {
		n.prefix = "[" + prefix + "] "
	} else {
		n.prefix = ""
	}
	n.init()
}

// SetFlags choose which flags the logger will display (after prefix and message
// level). See the package log for more informations on this.
func (n *Notepad) SetFlags(flags int) {
	n.flags = flags
	n.init()
}

// UseColor enabled or disables colored output.
func (n *Notepad) UseColor() {
	n.colored = true
}

// Feedback writes plainly to the outHandle while
// logging with the standard extra information (date, file, etc).
type Feedback struct {
	out *log.Logger
	log *log.Logger
}

func (fb *Feedback) Println(v ...interface{}) {
	fb.output(fmt.Sprintln(v...))
}

func (fb *Feedback) Printf(format string, v ...interface{}) {
	fb.output(fmt.Sprintf(format, v...))
}

func (fb *Feedback) Print(v ...interface{}) {
	fb.output(fmt.Sprint(v...))
}

func (fb *Feedback) output(s string) {
	if fb.out != nil {
		fb.out.Output(2, s)
	}
	if fb.log != nil {
		fb.log.Output(2, s)
	}
}
