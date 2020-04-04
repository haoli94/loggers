package mylogger

import (
	"fmt"
	"os"
	"path"
	"time"
)

type FileLogger struct {
	level       LogLevel
	filePath    string
	fileName    string
	maxFileSize int64
	fd          *os.File
	errfd       *os.File
}

func NewFileLogger(level, fp, fn string, maxFileSize int64) *FileLogger {
	levelInt, err := parseLogLevel(level)
	if err != nil {
		panic(err)
	}
	fl := &FileLogger{
		level:       levelInt,
		filePath:    fp,
		fileName:    fn,
		maxFileSize: maxFileSize,
	}
	err = fl.initFds()
	if err != nil {
		panic(err)
	}
	return fl
}

func (f *FileLogger) initFds() error {
	fullFileName := path.Join(f.filePath, f.fileName)
	fd, err := os.OpenFile(fullFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Open log file failed, err %v\n", err)
		return err
	}
	f.fd = fd

	errfd, err := os.OpenFile(fullFileName+".err", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Open error log file failed, err %v\n", err)
		return err
	}
	f.errfd = errfd
	return nil
}

func (f *FileLogger) Close() {
	f.errfd.Close()
	f.fd.Close()
}

func (f *FileLogger) checkSize(file *os.File) bool {
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("get file info failed, err %v\n", err)
		return false
	}
	return fileInfo.Size() >= f.maxFileSize
}
func (f *FileLogger) enable(level LogLevel) bool {
	if f.level < level {
		return true
	}
	return false
}

func (f *FileLogger) unParse(level LogLevel) string {
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

func (f *FileLogger) splitFile(file *os.File) (*os.File, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("get file info failed, err %v\n", err)
		return nil, err
	}
	logName := path.Join(f.filePath, fileInfo.Name())
	timeStamp := time.Now().Format("200601021504050000")
	newLogName := logName + timeStamp
	file.Close()
	os.Rename(logName, newLogName)
	fd, err := os.OpenFile(logName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open new log file failed, err %v\n", err)
		return nil, err
	}
	return fd, nil
}

func (f *FileLogger) log(t LogLevel, msg string, a ...interface{}) {
	if f.enable(t) {
		msg = fmt.Sprintf(msg, a...)
		funcName, fileName, line := getParams(3)
		now := time.Now()
		if f.checkSize(f.fd) {
			fd, err := f.splitFile(f.fd)
			if err != nil {
				fmt.Printf("split log file failed, err %v\n", err)
				return
			}
			f.fd = fd
		}
		fmt.Fprintf(f.fd, "[%s] [%s] [%s %s %d] %s \n", now.Format("2006-01-02 15:04:05"), f.unParse(t), fileName, funcName, line, msg)
		if t >= ERROR {
			if f.checkSize(f.fd) {
				errfd, err := f.splitFile(f.errfd)
				if err != nil {
					fmt.Printf("split err log file failed, err %v\n", err)
					return
				}
				f.errfd = errfd
			}
			fmt.Fprintf(f.errfd, "[%s] [%s] [%s %s %d] %s \n", now.Format("2006-01-02 15:04:05"), f.unParse(t), fileName, funcName, line, msg)
		}
	}
}

func (f *FileLogger) Debug(msg string, a ...interface{}) {
	f.log(DEBUG, msg, a...)
}

func (f *FileLogger) Trace(msg string, a ...interface{}) {
	f.log(TRACE, msg, a...)
}

func (f *FileLogger) Info(msg string, a ...interface{}) {
	f.log(INFO, msg, a...)
}

func (f *FileLogger) Warning(msg string, a ...interface{}) {
	f.log(WARNING, msg, a...)
}

func (f *FileLogger) Error(msg string, a ...interface{}) {
	f.log(ERROR, msg, a...)
}

func (f *FileLogger) Fatal(msg string, a ...interface{}) {
	f.log(FATAL, msg, a...)
}
