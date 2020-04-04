package mylogger

import (
	"errors"
	"fmt"
	"path"
	"runtime"
	"strings"
	"time"
)

type LogLevel uint16

const (
	UNKNOWN LogLevel = iota
	DEBUG
	TRACE
	INFO
	WARNING
	ERROR
	FATAL
)

type ConsoleLogger struct {
	level LogLevel
}

func parseLogLevel(s string) (LogLevel, error) {
	s = strings.ToLower(s)
	switch s {
	case "debug":
		return DEBUG, nil
	case "trace":
		return TRACE, nil
	case "info":
		return INFO, nil
	case "warning":
		return WARNING, nil
	case "error":
		return ERROR, nil
	case "fatal":
		return FATAL, nil
	default:
		return UNKNOWN, errors.New("Unknown log level.")
	}
}

func NewConsoleLogger(level string) *ConsoleLogger {
	levelInt, err := parseLogLevel(level)
	if err != nil {
		panic(err)
	}
	return &ConsoleLogger{
		level: levelInt,
	}
}

func getParams(skip int) (funcName, fileName string, line int) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		fmt.Println("runtime.Caller() failed")
		return
	}
	funcName = strings.Split(runtime.FuncForPC(pc).Name(), ".")[1]
	fileName = path.Base(file)
	return
}

func (c *ConsoleLogger) unParse(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case TRACE:
		return "TRACE"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return ""
	}
}

func (c *ConsoleLogger) enable(level LogLevel) bool {
	if c.level < level {
		return true
	}
	return false
}

func (c *ConsoleLogger) log(t LogLevel, msg string, a ...interface{}) {
	if c.enable(t) {
		msg = fmt.Sprintf(msg, a...)
		funcName, fileName, line := getParams(3)
		now := time.Now()
		fmt.Printf("[%s] [%s] [%s %s %d] %s \n", now.Format("2006-01-02 15:04:05"), c.unParse(t), fileName, funcName, line, msg)
	}
}

func (c *ConsoleLogger) Debug(msg string, a ...interface{}) {
	c.log(DEBUG, msg, a...)
}

func (c *ConsoleLogger) Trace(msg string, a ...interface{}) {
	c.log(TRACE, msg, a...)
}

func (c *ConsoleLogger) Info(msg string, a ...interface{}) {
	c.log(INFO, msg, a...)
}

func (c *ConsoleLogger) Warning(msg string, a ...interface{}) {
	c.log(WARNING, msg, a...)
}

func (c *ConsoleLogger) Error(msg string, a ...interface{}) {
	c.log(ERROR, msg, a...)
}

func (c *ConsoleLogger) Fatal(msg string, a ...interface{}) {
	c.log(FATAL, msg, a...)
}
