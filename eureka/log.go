package eureka

import (
    . "log"
)

type LogFunc func(format string, a ...interface{})

var log LogFunc

func (t LogFunc) Debugf(format string, a ...interface{}) {
    t(format, a...)
}

func init() {
    if log == nil {
        log = func(format string, a ...interface{}) {
            Printf(format, a...)
        }
    }
}
