// 日志工具类
package ezlog

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// 日志级别常量
// LogLevel Constants
const (
	LVL_DEBUG = 1 << iota
	LVL_INFO
	LVL_WARN
	LVL_ERROR
)

// Log
// Filename 文件路径
// Pattern  日期表达式（可选，默认无）
// Suffix   日志文件后缀（可选，默认"log"）
// LogLevel 日志级别（可选，默认"LVL_DEBUG"）
// curLogFile 当前日志文件
// buf for accumulating text to write
type Log struct {
	Filename   string
	Pattern    string
	Suffix     string
	LogLevel   int
	mu         sync.Mutex
	curLogFile *os.File
	buf        []byte
	isInited   bool
}

// New
func (l *Log) init() error {

	if !l.isInited {

		l.mu.Lock()
		defer l.mu.Unlock()

		if !l.isInited {

			// 准备日志文件，父文件夹
			_dir := filepath.Dir(l.Filename)
			exist, err := isPathExist(_dir)
			if err != nil {
				fmt.Printf("get dir error![%v]\n", err)
				return err
			}

			if !exist {
				err := os.MkdirAll(_dir, os.ModePerm)
				if err != nil {
					fmt.Printf("mkdir failed![%v]\n", err)
					return err
				}
			}

			// 准备默认参数
			if l.LogLevel == 0 {
				l.LogLevel = LVL_DEBUG
			}

			if l.Suffix == "" {
				l.Suffix = "log"
			}
			l.isInited = true
		}

	}

	return nil
}

// 判断文件夹是否存在
func isPathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 计算日志文件路径
func (l *Log) getLogPath(t *time.Time) string {
	var buffer bytes.Buffer
	buffer.WriteString(l.Filename)

	if l.Pattern != "" {
		buffer.WriteString(t.Format(l.Pattern))
	}

	buffer.WriteString(".")
	buffer.WriteString(l.Suffix)
	return buffer.String()
}

// 创建并打开新文件
func (l *Log) createAndOpenFile(filepath string) error {

	exist, err := isPathExist(filepath)
	if err != nil {
		fmt.Printf("func isPathExist error![%v]\n", err)
		return err
	} else {

		if !exist {

			file, err := os.Create(filepath)
			defer file.Close()
			if err != nil {
				fmt.Printf("create new log file error![%v]\n", err)
				return err
			}
		}
	}

	logFile, err := os.OpenFile(filepath, os.O_RDWR|os.O_APPEND, 0644)

	if err != nil {
		fmt.Printf("open new log file error![%v]\n", err)
		return err
	}

	l.curLogFile = logFile
	return nil
}

// 准备日志文件
func (l *Log) prepareLogFile(filepath string) error {

	if l.curLogFile != nil {

		if strings.Compare(l.curLogFile.Name(), filepath) != 0 {

			l.curLogFile.Close()
			err := l.createAndOpenFile(filepath)
			if err != nil {
				return err
			}

		}
	} else {

		err := l.createAndOpenFile(filepath)
		if err != nil {
			return err
		}

	}
	return nil
}

// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func itoa(buf *[]byte, i int, wid int) {

	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

// appendLevel
func appendLevel(buf *[]byte, level int) {

	var prefix string
	switch level {
	case LVL_DEBUG:
		prefix = "[Debug]"
	case LVL_INFO:
		prefix = "[Info]"
	case LVL_WARN:
		prefix = "[Warn]"
	case LVL_ERROR:
		prefix = "[Error]"
	}

	*buf = append(*buf, prefix...)
}

// 写日志
func (l *Log) writeLog(msg string, level int) error {
	err := l.init()
	if err != nil {
		fmt.Printf("init error![%v]\n", err)
		return err
	}

	if l.LogLevel <= level {

		l.mu.Lock()
		defer l.mu.Unlock()

		t := time.Now()
		err := l.prepareLogFile(l.getLogPath(&t))
		if err != nil {
			return err
		}

		l.buf = l.buf[:0]

		//format Header
		year, month, day := t.Date()
		itoa(&l.buf, year, 4)
		l.buf = append(l.buf, '/')
		itoa(&l.buf, int(month), 2)
		l.buf = append(l.buf, '/')
		itoa(&l.buf, day, 2)
		l.buf = append(l.buf, ' ')

		hour, min, sec := t.Clock()
		itoa(&l.buf, hour, 2)
		l.buf = append(l.buf, ':')
		itoa(&l.buf, min, 2)
		l.buf = append(l.buf, ':')
		itoa(&l.buf, sec, 2)

		l.buf = append(l.buf, '.')
		itoa(&l.buf, t.Nanosecond()/1e6, 3)
		l.buf = append(l.buf, ' ')

		// log level
		appendLevel(&l.buf, level)

		// log msg
		l.buf = append(l.buf, msg...)
		if len(msg) == 0 || msg[len(msg)-1] != '\n' {
			l.buf = append(l.buf, '\n')
		}
		_, err = l.curLogFile.Write(l.buf)
		if err != nil {
			fmt.Printf("write log error![%v]\n", err)
			return err
		}
	}
	return nil
}

// 业务日志，不受日志级别影响
// Arguments are handled in the manner of fmt.Printf.
func (l *Log) Printf(msg string, v ...interface{}) {
	l.writeLog(fmt.Sprintf(msg, v...), 999)
}

// 业务日志，不受日志级别影响
func (l *Log) Print(msg string) {
	l.writeLog(msg, 999)
}

// debug级别
func (l *Log) Debug(msg string) {
	l.writeLog(msg, LVL_DEBUG)
}

// info级别
func (l *Log) Info(msg string) {
	l.writeLog(msg, LVL_INFO)
}

// warn级别
func (l *Log) Warn(msg string) {
	l.writeLog(msg, LVL_WARN)
}

// error级别
func (l *Log) Error(msg string) {
	l.writeLog(msg, LVL_ERROR)
}
