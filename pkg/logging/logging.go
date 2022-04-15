package logging

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Logger struct {
	isActive bool
	writer   io.Writer
	prefix   string
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

func (l *Logger) SetPrefix(p string) {
	l.prefix = p
}

func (l *Logger) Printf(format string, a ...any) {
	if l.isActive {
		fmt.Fprint(l.writer, l.prefix)
		fmt.Fprintf(l.writer, format, a...)
	}
}

func (l *Logger) Write(s string) {
	for _, v := range strings.Split(s, "\n") {
		l.Printf("%s\n", v)
	}
}
