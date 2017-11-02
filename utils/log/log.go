package log

import (
	"fmt"
	"sync"
	"time"
)

var debug bool = true

type LogLevel string

const (
	DEBUG   = "D"
	INFO    = "I"
	WARNING = "W"
	ERROR   = "E"
	STATUE  = "S"
)

var statuLog string = ""
var lock sync.Mutex

func MyLogS(format string, a ...interface{}) {
	myLog(STATUE, format, a...)
}

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

func cleanStatus() {
	for i := 0; i < len(statuLog); i++ {
		fmt.Printf("\b")
	}
}

func Clean() {
	cleanStatus()
	fmt.Println()
}

func myLog(l LogLevel, format string, a ...interface{}) {
	if l == DEBUG && debug == false {
		return
	}
	newformat := string("[%v][%s]:") + format + string("\n")
	arg := make([]interface{}, 2, len(a)+3)
	arg[0] = time.Now().Format("2006-01-02 15:04:05.999")
	arg[1] = l
	if len(a) > 0 {
		arg = append(arg, a...)
	}
	var buf string
	if l == STATUE {
		buf = fmt.Sprintf(format, a...)
	} else {
		buf = fmt.Sprintf(newformat, arg...)
	}
	lock.Lock()
	defer lock.Unlock()
	cleanStatus()
	if l == STATUE {
		fmt.Print(buf)
		statuLog = buf
	} else {
		fmt.Print(buf)
		statuLog = ""
	}

}
