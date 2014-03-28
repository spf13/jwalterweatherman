// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package jwalterweatherman

import (
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "os"
)

// Level describes the chosen log level between
// debug and critical.
type Level int

//TRACE
//DEBUG
//INFO
//WARN
//ERROR
//CRITICAL
//FATAL

type JWWLevel struct {
    Handle io.Writer
    Level  int
    *log.Logger
}

type Feedback struct{}

// Log levels to control the logging output.
const (
    LevelDebug Level = iota
    LevelWarn
    LevelInfo
    LevelError
    LevelCritical
)

var (
    DEBUG          *log.Logger
    WARN           *log.Logger
    INFO           *log.Logger
    LOG            *log.Logger
    ERROR          *log.Logger
    CRITICAL       *log.Logger
    FEEDBACK       Feedback
    DebugHandle    io.Writer = os.Stdout
    WarnHandle     io.Writer = os.Stdout
    InfoHandle     io.Writer = os.Stdout
    ErrorHandle    io.Writer = os.Stdout
    CriticalHandle io.Writer = os.Stdout
    logLevel       Level     = LevelWarn // 1
    outLevel       Level     = LevelInfo // 2
    LogHandle      io.Writer = ioutil.Discard
    OutHandle      io.Writer = os.Stdout
    BothHandle     io.Writer = io.MultiWriter(LogHandle, OutHandle)
)

func init() {
    SetOutLevel(LevelInfo)
}

func Initialize() {
    initWriters()

    DEBUG = log.New(DebugHandle,
        "DEBUG: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    INFO = log.New(InfoHandle,
        "INFO:  ",
        log.Ldate|log.Ltime|log.Lshortfile)

    LOG = log.New(LogHandle,
        "LOG:   ",
        log.Ldate|log.Ltime|log.Lshortfile)

    WARN = log.New(WarnHandle,
        "WARN:  ",
        log.Ldate|log.Ltime|log.Lshortfile)

    ERROR = log.New(ErrorHandle,
        "ERROR: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    CRITICAL = log.New(CriticalHandle,
        "CRITICAL: ",
        log.Ldate|log.Ltime|log.Lshortfile)
}

// Level returns the current log level.
func LogLevel() Level {
    return logLevel
}

func OutLevel() Level {
    return outLevel
}

func levelCheck(level Level) Level {
    switch {
    case level <= LevelDebug:
        return LevelDebug
    case level >= LevelCritical:
        return LevelCritical
    default:
        return level
    }
}

// SetLevel switches to a new log level.
func SetLogLevel(level Level) {
    logLevel = levelCheck(level)
    Initialize()
}

func SetLogFile(path string) {
    file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
    fmt.Println("Logging to", file.Name())
    if err != nil {
        CRITICAL.Println("Failed to open log file:", path, err)
        os.Exit(-1)
    }

    LogHandle = file
    Initialize()
}

func UseTempLogFile(prefix string) {
    file, err := ioutil.TempFile(os.TempDir(), prefix)
    if err != nil {
        CRITICAL.Println(err)
    }

    fmt.Println("Logging to", file.Name())

    LogHandle = file
    Initialize()
}

func DiscardLogging() {
    LogHandle = ioutil.Discard
    Initialize()
}

func SetOutLevel(level Level) {
    outLevel = levelCheck(level) // 1
    Initialize()
}

// Don't use if you have manually set the Handles of the different levels as it will overwrite them.
func initWriters() {
    BothHandle = io.MultiWriter(LogHandle, OutHandle)
    //DEBUG
    if LevelDebug < outLevel && LevelDebug < logLevel {
        DebugHandle = ioutil.Discard
    } else if LevelDebug >= outLevel && LevelDebug >= logLevel {
        DebugHandle = BothHandle
    } else if LevelDebug >= outLevel && LevelDebug < logLevel {
        DebugHandle = OutHandle
    } else {
        DebugHandle = LogHandle
    }

    //WARN
    if LevelWarn < outLevel && LevelWarn < logLevel {
        WarnHandle = ioutil.Discard
    } else if LevelWarn >= outLevel && LevelWarn >= logLevel {
        WarnHandle = BothHandle
    } else if LevelWarn >= outLevel && LevelWarn < logLevel {
        WarnHandle = OutHandle
    } else {
        WarnHandle = LogHandle
    }

    //INFO
    if LevelInfo < outLevel && LevelInfo < logLevel {
        InfoHandle = ioutil.Discard
    } else if LevelInfo >= outLevel && LevelInfo >= logLevel {
        InfoHandle = BothHandle
    } else if LevelInfo >= outLevel && LevelInfo < logLevel {
        InfoHandle = OutHandle
    } else {
        InfoHandle = LogHandle
    }

    //ERROR
    if LevelError < outLevel && LevelError < logLevel {
        ErrorHandle = ioutil.Discard
    } else if LevelError >= outLevel && LevelError >= logLevel {
        ErrorHandle = BothHandle
    } else if LevelError >= outLevel && LevelError < logLevel {
        ErrorHandle = OutHandle
    } else {
        ErrorHandle = LogHandle
    }

    //CRITICAL
    if LevelCritical < outLevel && LevelCritical < logLevel {
        CriticalHandle = ioutil.Discard
    } else if LevelCritical >= outLevel && LevelCritical >= logLevel {
        CriticalHandle = BothHandle
    } else if LevelCritical >= outLevel && LevelCritical < logLevel {
        CriticalHandle = OutHandle
    } else {
        CriticalHandle = LogHandle
    }
}

func (fb *Feedback) Println(v ...interface{}) {
    fmt.Println(v...)
    LOG.Println(v...)
}

func (fb *Feedback) Printf(format string, v ...interface{}) {
    fmt.Printf(format, v...)
    LOG.Printf(format, v...)
}
