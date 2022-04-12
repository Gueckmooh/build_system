package logging

import (
	"fmt"
	"io"
	"os"
)

type Logger struct {
	isActive bool
	writer   io.Writer
}

func NewLogger(writer io.Writer) *Logger {
	return &Logger{
		isActive: false,
		writer:   writer,
	}
}

var (
	Debug *Logger
	Log   *Logger
)

func init() {
	Debug = NewLogger(os.Stdout)
	Log = NewLogger(os.Stdout)
}

func SetDebugLogging(v bool) {
	Debug.isActive = true
}

func SetVerboseLogging(v bool) {
	Log.isActive = true
}

func (l *Logger) Printf(format string, a ...any) {
	if l.isActive {
		fmt.Fprintf(l.writer, format, a...)
	}
}
