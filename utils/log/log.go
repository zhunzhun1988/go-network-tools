package log

import (
	"fmt"
	"time"
)

var debug bool = true

type LogLevel string

const (
	DEBUG   = "D"
	INFO    = "I"
	WARNING = "W"
	ERROR   = "E"
)

func MyLogI(format string, a ...interface{}) {
	myLog(INFO, format, a...)
}
func MyLogW(format string, a ...interface{}) {
	myLog(WARNING, format, a...)
}
func MyLogD(format string, a ...interface{}) {
	myLog(DEBUG, format, a...)
}
func MyLogE(format string, a ...interface{}) {
	myLog(ERROR, format, a...)
}

func myLog(l LogLevel, format string, a ...interface{}) {
	if l == DEBUG && debug == false {
		return
	}
	format = string("[%v][%s]:") + format + string("\n")
	arg := make([]interface{}, 2, len(a)+3)
	arg[0] = time.Now().Format("2006-01-02 15:04:05.999")
	arg[1] = l
	if len(a) > 0 {
		arg = append(arg, a...)
	}
	fmt.Printf(format, arg...)
}
