package main

import (
	ezlog "github.com/sialot/ezlog"
)

func main() {

	config := log.Config{Filename: "e:/goTest/newLogger"}
	logger := log.New(config)

	logger.Print("this is a test log.")
	logger.Printf("this is a test log. %d", 123)
	logger.Debug("debug msg")
	logger.Info("info msg")
	logger.Warn("warn msg")
	logger.Error("error msg")

}
