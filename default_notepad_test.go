// Copyright © 2016 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package jwalterweatherman

import (
	"bytes"
	"io/ioutil"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestThresholds(t *testing.T) {
	origStdoutThresh := defaultNotepad.stdoutThreshold
	origLogThresh := defaultNotepad.logThreshold
	defer func() {
		defaultNotepad.stdoutThreshold = origStdoutThresh
		defaultNotepad.logThreshold = origLogThresh
		defaultNotepad.init()
		reloadDefaultNotepad()
	}()

	SetStdoutThreshold(LevelError)
	require.Equal(t, StdoutThreshold(), LevelError)
	SetLogThreshold(LevelCritical)
	require.Equal(t, LogThreshold(), LevelCritical)
	require.NotEqual(t, StdoutThreshold(), LevelCritical)
	SetStdoutThreshold(LevelWarn)
	require.Equal(t, StdoutThreshold(), LevelWarn)
}

func TestDefaultLogging(t *testing.T) {
	var outputBuf, errBuf, logBuf bytes.Buffer

	origOut := defaultNotepad.outHandle
	origErr := defaultNotepad.errHandle
	origLog := defaultNotepad.logHandle
	origStdoutThresh := defaultNotepad.stdoutThreshold
	origLogThresh := defaultNotepad.logThreshold
	defer func() {
		defaultNotepad.outHandle = origOut
		defaultNotepad.errHandle = origErr
		defaultNotepad.logHandle = origLog
		defaultNotepad.stdoutThreshold = origStdoutThresh
		defaultNotepad.logThreshold = origLogThresh
		defaultNotepad.init()
		reloadDefaultNotepad()
	}()

	defaultNotepad.logHandle = &logBuf
	defaultNotepad.outHandle = &outputBuf
	defaultNotepad.errHandle = &errBuf

	SetLogThreshold(LevelWarn)
	SetStdoutThreshold(LevelError)

	FATAL.Println("fatal err")
	CRITICAL.Println("critical err")
	ERROR.Println("an error")
	WARN.Println("a warning")
	INFO.Println("information")
	DEBUG.Println("debugging info")
	TRACE.Println("trace")

	require.Contains(t, logBuf.String(), "fatal err")
	require.Contains(t, logBuf.String(), "critical err")
	require.Contains(t, logBuf.String(), "an error")
	require.Contains(t, logBuf.String(), "a warning")
	require.NotContains(t, logBuf.String(), "information")
	require.NotContains(t, logBuf.String(), "debugging info")
	require.NotContains(t, logBuf.String(), "trace")

	require.Contains(t, errBuf.String(), "fatal err")
	require.Contains(t, errBuf.String(), "critical err")
	require.Contains(t, errBuf.String(), "an error")
	require.NotContains(t, errBuf.String(), "a warning")
	require.NotContains(t, errBuf.String(), "information")
	require.NotContains(t, errBuf.String(), "debugging info")
	require.NotContains(t, errBuf.String(), "trace")

	require.NotContains(t, outputBuf.String(), "fatal err")
	require.NotContains(t, outputBuf.String(), "critical err")
	require.NotContains(t, outputBuf.String(), "an error")
	require.NotContains(t, outputBuf.String(), "a warning")
	require.NotContains(t, outputBuf.String(), "information")
	require.NotContains(t, outputBuf.String(), "debugging info")
	require.NotContains(t, outputBuf.String(), "trace")
}

func TestLogCounter(t *testing.T) {
	assert := require.New(t)

	origOut := defaultNotepad.outHandle
	origErr := defaultNotepad.errHandle
	origLog := defaultNotepad.logHandle
	origStdoutThresh := defaultNotepad.stdoutThreshold
	origLogThresh := defaultNotepad.logThreshold
	defer func() {
		defaultNotepad.outHandle = origOut
		defaultNotepad.errHandle = origErr
		defaultNotepad.logHandle = origLog
		defaultNotepad.stdoutThreshold = origStdoutThresh
		defaultNotepad.logThreshold = origLogThresh
		defaultNotepad.init()
		reloadDefaultNotepad()
	}()

	defaultNotepad.logHandle = ioutil.Discard
	defaultNotepad.outHandle = ioutil.Discard
	defaultNotepad.errHandle = ioutil.Discard

	errorCounter := &Counter{}

	SetLogThreshold(LevelTrace)
	SetStdoutThreshold(LevelTrace)
	SetLogListeners(LogCounter(errorCounter, LevelError))

	FATAL.Println("fatal err")
	CRITICAL.Println("critical err")
	WARN.Println("a warning")
	WARN.Println("another warning")
	INFO.Println("information")
	DEBUG.Println("debugging info")
	TRACE.Println("trace")

	wg := &sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				ERROR.Println("error", j)
				// check for data races
				assert.True(errorCounter.Count() > uint64(j))
			}
		}()

	}

	wg.Wait()

	assert.Equal(uint64(102), errorCounter.Count())
}

func TestDefaultErrorLoggersStderr(t *testing.T) {
	var outputBuf, errBuf bytes.Buffer

	origOut := defaultNotepad.outHandle
	origErr := defaultNotepad.errHandle
	origStdoutThresh := defaultNotepad.stdoutThreshold
	defer func() {
		defaultNotepad.outHandle = origOut
		defaultNotepad.errHandle = origErr
		defaultNotepad.stdoutThreshold = origStdoutThresh
		defaultNotepad.init()
		reloadDefaultNotepad()
	}()

	SetStdoutOutput(&outputBuf)
	SetStderrOutput(&errBuf)
	SetStdoutThreshold(LevelInfo)

	INFO.Println("info message")
	WARN.Println("warn message")
	ERROR.Println("error message")
	CRITICAL.Println("critical message")
	FATAL.Println("fatal message")

	// Info and Warn must route to stdout
	require.Contains(t, outputBuf.String(), "info message")
	require.Contains(t, outputBuf.String(), "warn message")
	require.NotContains(t, outputBuf.String(), "error message")
	require.NotContains(t, outputBuf.String(), "critical message")
	require.NotContains(t, outputBuf.String(), "fatal message")

	// Error, Critical, and Fatal must route to stderr
	require.NotContains(t, errBuf.String(), "info message")
	require.NotContains(t, errBuf.String(), "warn message")
	require.Contains(t, errBuf.String(), "error message")
	require.Contains(t, errBuf.String(), "critical message")
	require.Contains(t, errBuf.String(), "fatal message")
}

func TestDefaultHandles(t *testing.T) {
	require.NotNil(t, GetStdoutOutput())
	require.NotNil(t, GetStderrOutput())
}
