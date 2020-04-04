package main

import (
	"errors"
	"loggers/mylogger"
	"time"
)

var myLogger mylogger.Logger

func main() {
	//myLogger = mylogger.NewFileLogger("info", "./", "hao.log", 256)
	myLogger = mylogger.NewConsoleLogger("info")
	for {
		myLogger.Debug("This is Debug")
		myLogger.Trace("This is Trace")
		myLogger.Info("This is Info")
		myLogger.Warning("This is Warning")
		myLogger.Error("This is Error", errors.New("No Error."), errors.New("hello world!"))
		myLogger.Fatal("This is Fatal")
		time.Sleep(1 * time.Second)
	}
}
