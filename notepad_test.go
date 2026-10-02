// Copyright © 2016 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package jwalterweatherman

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNotepad(t *testing.T) {
	var logHandle, outHandle bytes.Buffer

	errorCounter := &Counter{}

	n := NewNotepad(LevelCritical, LevelError, &outHandle, &logHandle, "TestNotePad", 0, LogCounter(errorCounter, LevelError))

	require.Equal(t, LevelCritical, n.GetStdoutThreshold())
	require.Equal(t, LevelError, n.GetLogThreshold())

	n.DEBUG.Println("Some debug")
	n.ERROR.Println("Some error")
	n.CRITICAL.Println("Some critical error")

	require.Contains(t, logHandle.String(), "[TestNotePad] ERROR Some error")
	require.NotContains(t, logHandle.String(), "Some debug")
	require.NotContains(t, outHandle.String(), "Some error")
	require.Contains(t, outHandle.String(), "Some critical error")

	// 1 error + 1 critical
	require.Equal(t, errorCounter.Count(), uint64(2))
}

func TestNotepadLogListener(t *testing.T) {
	assert := require.New(t)

	var errorBuff, infoBuff bytes.Buffer

	errorCapture := func(t Threshold) io.Writer {
		if t != LevelError {
			// Only interested in ERROR
			return nil
		}

		return &errorBuff
	}

	infoCapture := func(t Threshold) io.Writer {
		if t != LevelInfo {
			return nil
		}

		return &infoBuff
	}

	n := NewNotepad(LevelCritical, LevelError, ioutil.Discard, ioutil.Discard, "TestNotePad", 0, infoCapture, errorCapture)

	n.DEBUG.Println("Some debug")
	n.INFO.Println("Some info")
	n.INFO.Println("Some more info")
	n.ERROR.Println("Some error")
	n.CRITICAL.Println("Some critical error")
	n.ERROR.Println("Some more error")

	assert.Equal(`[TestNotePad] ERROR Some error
[TestNotePad] ERROR Some more error
`, errorBuff.String())
	assert.Equal(`[TestNotePad] INFO Some info
[TestNotePad] INFO Some more info
`, infoBuff.String())

}

func TestThresholdString(t *testing.T) {
	require.Equal(t, LevelError.String(), "ERROR")
	require.Equal(t, LevelTrace.String(), "TRACE")
}

func TestCustomThresholdPrefix(t *testing.T) {
	var logHandle, outHandle bytes.Buffer

	n := NewNotepad(LevelTrace, LevelTrace, &outHandle, &logHandle, "TestCustom", 0)

	require.Equal(t, "INFO", n.GetThresholdPrefix(LevelInfo))
	require.Equal(t, "WARN", n.GetThresholdPrefix(LevelWarn))

	n.SetThresholdPrefix(LevelInfo, "INF")
	n.SetThresholdPrefix(LevelWarn, "WRN")

	require.Equal(t, "INF", n.GetThresholdPrefix(LevelInfo))
	require.Equal(t, "WRN", n.GetThresholdPrefix(LevelWarn))
	require.Equal(t, "DEBUG", n.GetThresholdPrefix(LevelDebug))

	n.INFO.Println("info msg")
	n.WARN.Println("warn msg")
	n.DEBUG.Println("debug msg")

	require.Contains(t, outHandle.String(), "[TestCustom] INF info msg")
	require.Contains(t, outHandle.String(), "[TestCustom] WRN warn msg")
	require.Contains(t, outHandle.String(), "[TestCustom] DEBUG debug msg")
}

func TestSetThresholdPrefixes(t *testing.T) {
	var logHandle, outHandle bytes.Buffer

	n := NewNotepad(LevelTrace, LevelTrace, &outHandle, &logHandle, "", 0)

	n.SetThresholdPrefixes(map[Threshold]string{
		LevelTrace:    "TRAC",
		LevelDebug:    "DEBG",
		LevelInfo:     "INFO",
		LevelWarn:     "WARN",
		LevelError:    "ERRO",
		LevelCritical: "CRIT",
		LevelFatal:    "FATL",
	})

	n.TRACE.Println("trace msg")
	n.DEBUG.Println("debug msg")
	n.ERROR.Println("error msg")
	n.CRITICAL.Println("crit msg")
	n.FATAL.Println("fatal msg")

	require.Contains(t, outHandle.String(), "TRAC trace msg")
	require.Contains(t, outHandle.String(), "DEBG debug msg")
	require.Contains(t, outHandle.String(), "ERRO error msg")
	require.Contains(t, outHandle.String(), "CRIT crit msg")
	require.Contains(t, outHandle.String(), "FATL fatal msg")
}

func TestEmptyThresholdPrefix(t *testing.T) {
	var logHandle, outHandle bytes.Buffer

	n := NewNotepad(LevelTrace, LevelTrace, &outHandle, &logHandle, "", 0)

	n.SetThresholdPrefix(LevelInfo, "")
	n.INFO.Println("bare message")

	require.Equal(t, "bare message\n", outHandle.String())
}

func BenchmarkLogPrintOnlyToCounter(b *testing.B) {
	var logHandle, outHandle bytes.Buffer
	n := NewNotepad(LevelCritical, LevelCritical, &outHandle, &logHandle, "TestNotePad", 0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n.INFO.Print("Test")
	}
}
