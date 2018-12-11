
# Ezlog 

Ezlog 是一个Go语言的日志输出工具，支持日志文件分割和日志缓冲功能。

### 初始化

Ezlog支持的初始化参数如下：

- Filename      文件路径
- Pattern       可选，日期表达式（默认:""）
- Suffix        可选，日志文件后缀（默认:"log"）
- LogLevel      可选，日志级别（默认:"LVL_DEBUG"）
- BufferSize    可选，缓存容量(默认:0，代表不启用缓存)

在指定了 Filename、Pattern、Suffix时，文件输出路径会按：Filename + Pattern + "." + Suffix 拼接

#### 最简初始化
日志的log级别默认为debug，输出日志位置：/var/log/demo.log

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
	logger := ezlog.Log{Filename: "/var/log/demo", LogLevel: ezlog.LVL_INFO}
```

#### 指定日志文件后缀

指定文件后缀为“txt”,输出日志位置：/var/log/demo.txt

```go
	logger := ezlog.Log{Filename: "/var/log/demo", Suffix:   "txt"}
```

#### 日志分割

Pattern 参数支持时间表达式。例如按秒分割文件"-2006-01-02_150405"，如果当前时间为："2018-12-11 11：28：00"，输出日志位置就是 "/var/log/demo-2018-12-11_112800.log"

```go
	logger := ezlog.Log{Filename: "/var/log/demo", Pattern:  "-2006-01-02_150405"}	
```

#### 日志缓冲

在初始化参数设置BufferSize，可以设置缓冲区大小，ezlog会在日志信息超过缓冲大小时才将缓冲区的数据写入文件，提高日志写入性能。为了保证没有数据丢失，请在应用中使用```func (l *Log) Flush() error ```方法手动清空缓冲区数据。

```go
	logger := ezlog.Log{Filename: "/var/log/demo",BufferSize:  1024}
	logger.Print("this is a test log.")
	...
	logger.Flush()
	
```
**在设置BufferSize之后，ezlog还会默认每200ms自动清空缓冲区并将缓冲区的数据写入文件**。可以使用```func (l *Log) SetFlushDuration(duration int) error ```方法指定时间间隔，参数单位为毫秒（ms）
```go
	logger := ezlog.Log{Filename: "/var/log/demo",BufferSize:  1024}
	logger.SetFlushDuration(100)
	...
```
如果不想使用自动缓冲清空功能，可以使用```func (l *Log) DisableAutoFlush() error ```方法。
```go
	logger := ezlog.Log{Filename: "/var/log/demo",BufferSize:  1024}
	logger.DisableAutoFlush()
        logger.Print("this is a test log.")
	...
	logger.Flush()
```

### 方法

#### Print(msg string) 
输出指定字符串，该方法会始终打印日志，不受日志级别设置的影响
```go
	logger.Print("this is a test log.")
```
#### Printf(msg string) 
格式化输出，该方法会始终打印日志，不受日志级别设置的影响
```go
	logger.Printf("this is a test log. %d", 123)
```
#### Debug(msg string) 
输出指定字符串，Debug 级别
```go
	logger.Debug("this is a debug log.")
```
#### Info(msg string) 
输出指定字符串，Info 级别
```go
	logger.Info("this is a info log.")
```
#### Warn(msg string) 
输出指定字符串，Warn 级别
```go
	logger.Warn("this is a warn log.")
```
#### Error(msg string) 
输出指定字符串，Error 级别
```go
	logger.Error("this is a error log.")
```
### 示例

```go
	package main
	
	import (
		"github.com/sialot/ezlog"
	)
	
	func main() {
	
		logger := ezlog.Log{
			Filename:   "/var/log/demo",
			Pattern:    "-2006-01-02_150405",
			Suffix:     "txt",
			LogLevel:   ezlog.LVL_DEBUG,
			BufferSize: 1024}
	
		logger.Print("this is a test log.")
		logger.Printf("this is a test log. %d", 123)
		logger.Debug("debug msg")
		logger.Info("info msg")
		logger.Warn("warn msg")
		logger.Error("error msg")
	
	}
```

