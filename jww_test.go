// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package jwalterweatherman

import (
    "bytes"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestLevels(t *testing.T) {
    SetOutLevel(LevelError)
    assert.Equal(t, OutLevel(), LevelError)
    SetLogLevel(LevelCritical)
    assert.Equal(t, LogLevel(), LevelCritical)
    assert.NotEqual(t, OutLevel(), LevelCritical)
    SetOutLevel(LevelWarn)
    assert.Equal(t, OutLevel(), LevelWarn)
}

func TestDefaultLogging(t *testing.T) {
    outputBuf := new(bytes.Buffer)
    logBuf := new(bytes.Buffer)
    LogHandle = logBuf
    OutHandle = outputBuf

    SetOutLevel(LevelInfo)
    SetLogLevel(LevelWarn)

    CRITICAL.Println("critical err")
    ERROR.Println("an error")
    WARN.Println("a warning")
    INFO.Println("information")
    DEBUG.Println("debugging info")

    assert.Contains(t, logBuf.String(), "critical err")
    assert.Contains(t, logBuf.String(), "an error")
    assert.Contains(t, logBuf.String(), "information")
    assert.Contains(t, logBuf.String(), "a warning")
    assert.NotContains(t, logBuf.String(), "debugging info")
    assert.Contains(t, outputBuf.String(), "critical err")
    assert.Contains(t, outputBuf.String(), "an error")
    assert.Contains(t, outputBuf.String(), "information")
    assert.NotContains(t, outputBuf.String(), "a warning")
    assert.NotContains(t, outputBuf.String(), "debugging info")
}
