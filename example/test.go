package main

import (
	"github.com/sialot/ezlog"
)

func main() {

	logger := &ezlog.Log{
		Filename: "/var/log/demo",
		Pattern:  "-2006-01-02_150405",
		Suffix:   "txt",
		LogLevel: ezlog.LVL_DEBUG}

	logger.Print("this is a test log.")
	logger.Printf("this is a test log. %d", 123)
	logger.Debug("debug msg")
	logger.Info("info msg")
	logger.Warn("warn msg")
	logger.Error("error msg")

}
