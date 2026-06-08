package utils

import "fmt"

type Logger struct {
	IsVerbose bool
}

func (l *Logger) Log(format string, args ...any) {
	if !l.IsVerbose {
		return
	}
	fmt.Printf(format+"\n", args...)
}
