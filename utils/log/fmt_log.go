package log

import (
	"fmt"
)

type FMTLogger struct {
}

func (l *FMTLogger) Info(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (l *FMTLogger) Debug(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (l *FMTLogger) Warn(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (l *FMTLogger) Error(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
