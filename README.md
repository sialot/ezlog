# ezlog
go语言的文件日志输出工具

## 初始化参数
- Filename 文件路径
- Pattern  日期表达式（可选，默认无）
- Suffix   日志文件后缀（可选，默认"log"）
- LogLevel 日志级别（可选，默认"LVL_DEBUG"）

文件输出路径按：Filename + Pattern + "." + Suffix 拼接

## 示例：

最小参数：

	package main
	
	import (
		"github.com/sialot/ezlog"
	)
	
	func main() {
	
		logger := ezlog.Log{
			Filename: "/var/log/demo"}
	
		logger.Print("this is a test log.")
		logger.Printf("this is a test log. %d", 123)
		logger.Debug("debug msg")
		logger.Info("info msg")
		logger.Warn("warn msg")
		logger.Error("error msg")
	
	}



所有参数：

	package main
	
	import (
		"github.com/sialot/ezlog"
	)
	
	func main() {
	
		logger := ezlog.Log{
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
