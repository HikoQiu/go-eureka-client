package eureka

import (
    . "log"
)

const (
    LevelDebug = 1
    LevelError = 3
)

type LogFunc func(level int, format string, a ...interface{})

var log LogFunc

func SetLogger(logFunc LogFunc) {
    log = logFunc
}

func (t LogFunc) Debugf(format string, a ...interface{}) {
    t(LevelDebug, format, a...)
}

func (t LogFunc) Errorf(format string, a ...interface{}) {
    t(LevelError, format, a...)
}

func init() {
    if log == nil {
        log = func(level int, format string, a ...interface{}) {
            switch level {
            case LevelDebug:
                format = "[debug] " + format
            case LevelError:
                format = "[error] " + format
            }
            Printf(format, a...)
        }
    }
}
