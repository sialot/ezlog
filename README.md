
# Ezlog 

Ezlog 是一个支持日志文件分割的日志包。


### 初始化参数

- Filename      文件路径
- Pattern       日期表达式（可选，默认无）
- Suffix        日志文件后缀（可选，默认"log"）
- LogLevel      日志级别（可选，默认"LVL_DEBUG"）
- BufferSize    缓存容量(可选，默认"0"，不启用缓存)

文件输出路径，按：Filename + Pattern + "." + Suffix 拼接

#### 最小初始化参数
日志的log级别默认为debug，输出日志位置：/var/log/demo-20181211_112800.txt

```go
	logger := ezlog.Log{Filename: "/var/log/demo"}	
```

#### 指定日志级别

日志级别常量：
- LVL_DEBUG
- LVL_INFO
- LVL_WARN
- LVL_ERROR

指定日志级别为info

```go

	logger := ezlog.Log{
			Filename: "/var/log/demo",			
			LogLevel: ezlog.LVL_INFO}
	
```

#### 指定日志文件后缀

指定文件后缀为“txt”,输出日志位置：/var/log/demo.txt

```go

	logger := ezlog.Log{
			Filename: "/var/log/demo",			
			Suffix:   "txt"}
	
```

#### 日志分割

Pattern 参数支持时间表达式，按秒分割文件

```go

	logger := ezlog.Log{
			Filename: "/var/log/demo",
			Pattern:  "-2006-01-02_150405",
			Suffix:   "txt",
			LogLevel: ezlog.LVL_DEBUG}
	
```

当前时间为：2018-12-11 11：28：00。输出日志位置：/var/log/demo-2018-12-11_112800.txt

#### 方法
```go
	func (l *Log) DisableAutoFlush() error 
	func (l *Log) SetFlushDuration(duration int) 
	func (l *Log) Flush() error 

	func (l *Log) Print(msg string) 
	func (l *Log) Printf(msg string, v ...interface{}) 
	func (l *Log) Debug(msg string) 
	func (l *Log) Info(msg string) 
	func (l *Log) Warn(msg string) 
	func (l *Log) Error(msg string) 
```
#### 示例

所有参数：

```go
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
```

