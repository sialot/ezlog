# ezlog
go语言的文件日志输出工具

## 初始化参数
- filename 文件路径
- pattern  日期表达式（可选，默认无）
- suffix   日志文件后缀（可选，默认"log"）
- logLevel 日志级别（可选，默认"LVL_DEBUG"）

文件输出路径按：filename + pattern + "." + suffix 拼接

## 示例：

默认参数：

	config := log.Config{Filename: "/var/log/demo/newLogger"}
	logger := log.New(config)

	logger.Print("this is a test log.")
	logger.Printf("this is a test log. %d", 123)
	logger.Debug("debug msg")
	logger.Info("info msg")
	logger.Warn("warn msg")
	logger.Error("error msg")

所有参数：

	config := log.Config{Filename: "/var/log/demo/newLogger", Pattern: "2006-01-02_150405", suffix:"txt"}
	logger.Print("this is a test log.")
	logger.Printf("this is a test log. %d", 123)
	logger.Debug("debug msg")
	logger.Info("info msg")
	logger.Warn("warn msg")
	logger.Error("error msg")
