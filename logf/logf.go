package logf

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"path"
	"runtime"
	"strconv"
	"time"
)

type LEVEL int
type Option int

const (
	DEBUG LEVEL = iota
	INFO
	WARN
	ERROR
	FATAL
)

const (
	NONE Option = 1 << iota
	TASK
)

type Options func(*Log)

type Log struct {
	option     Option
	skip       int
	timeFormat string
}

func Parser(option Option) Options {
	return func(l *Log) {
		l.option = option
	}
}

func Logger(opts ...Options) (c *Log) {
	c = &Log{timeFormat: time.Now().Format("2006/01/02 15:04:05"), skip: 3}
	for _, opt := range opts {
		opt(c)
	}
	return
}

func (c *Log) Skip(skip int) *Log {
	c.skip = skip
	return c
}

func (c *Log) caller() (filename string) {
	_, file, line, ok := runtime.Caller(c.skip)
	if !ok {
		file = "unknown"
		line = 0
	}
	filename = path.Base(file) + ":" + strconv.Itoa(line)
	return
}

func (c *Log) level(expr LEVEL) (level string) {
	switch expr {
	case INFO:
		level = color.CyanString("[INFO]")
	case WARN:
		level = color.YellowString("[WARN]")
	case ERROR:
		level = color.RedString("[ERROR]")
	}
	if c.option&TASK != 0 {
		level += " [TASK]"
	}
	return
}

func (c *Log) println(level string, value any) {
	if level == "" {
		_, _ = fmt.Fprintln(os.Stdout, c.timeFormat, c.caller(), value)
		return
	}
	_, _ = fmt.Fprintln(os.Stdout, c.timeFormat, c.caller(), level, value)
}

func (c *Log) Println(value any) {
	c.println("", value)
}

func (c *Log) Printf(format string, value ...any) {
	c.println("", fmt.Sprintf(format, value...))
}

func (c *Log) FatalLn(value any) {
	c.println("", value)
	os.Exit(1)
}

func (c *Log) Fatalf(format string, value ...any) {
	c.println("", fmt.Sprintf(format, value...))
	os.Exit(1)
}

func (c *Log) Info(value any) {
	c.println(c.level(INFO), value)
}

func (c *Log) InfoF(format string, value ...any) {
	c.println(c.level(INFO), fmt.Sprintf(format, value...))
}

func (c *Log) Error(value any) {
	c.println(c.level(ERROR), value)
}

func (c *Log) ErrorF(format string, value ...any) {
	c.println(c.level(ERROR), fmt.Sprintf(format, value...))
}

func (c *Log) Warn(err error) {
	c.println(c.level(WARN), err.Error())
}

func (c *Log) WarnIf(err error) {
	if err != nil {
		c.println(c.level(WARN), err.Error())
	}
}

func (c *Log) ErrorIf(err error) {
	if err != nil {
		c.println(c.level(ERROR), err.Error())
	}
}

func Sprintf(format string, value ...any) string {
	return fmt.Sprintf(format, value...)
}

func Error(format string, value ...any) error {
	return fmt.Errorf(format, value...)
}
