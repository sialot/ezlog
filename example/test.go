package main

import (
	ezlog "github.com/sialot/ezlog"
)

func main() {

	logger := ezlog.New("/var/log/ezlog-", "2006-01-02", "log", ezlog.LVL_DEBUG)

	logger.Print("this is a test log.")
	logger.Printf("this is a test log. %d", 123)
	logger.Debug("debug msg")
	logger.Info("info msg")
	logger.Warn("warn msg")
	logger.Error("error msg")

}