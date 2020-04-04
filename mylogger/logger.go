package mylogger

type Logger interface {
	Debug(msg string, a ...interface{})
	Trace(msg string, a ...interface{})
	Info(msg string, a ...interface{})
	Warning(msg string, a ...interface{})
	Error(msg string, a ...interface{})
	Fatal(msg string, a ...interface{})
}
