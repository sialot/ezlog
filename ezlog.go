// 日志工具类
package ezlog

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// 日志级别常量
const (
	LVL_DEBUG = 1 << iota
	LVL_INFO
	LVL_WARN
	LVL_ERROR
)

// 日志对象类型
type Log struct {
	mu         sync.Mutex
	filename   string
	pattern    string
	suffix     string
	logLevel   int
	curLogFile *os.File
	buf        []byte // for accumulating text to write
}

// filename 文件路径
// pattern  日期表达式
// suffix   日志后缀
// logLevel 日志级别
func New(filename string, pattern string, suffix string, logLevel int) *Log {

	// 准备日志文件，父文件夹
	_dir := filepath.Dir(filename)
	exist, err := isPathExist(_dir)
	if err != nil {
		fmt.Printf("get dir error![%v]\n", err)
		return nil
	}

	if !exist {
		err := os.MkdirAll(_dir, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir failed![%v]\n", err)
			return nil
		}
	}

	return &Log{filename: filename, pattern: pattern, suffix: suffix, logLevel: logLevel}
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

// 计算按日期格式的文件路径
func (l *Log) getLogPath(t time.Time) string {
	var buffer bytes.Buffer
	buffer.WriteString(l.filename)
	buffer.WriteString(t.Format(l.pattern))
	buffer.WriteString(".")
	buffer.WriteString(l.suffix)
	return buffer.String()
}

// 创建新文件
func (l *Log) createNewFile(filepath string) {

	exist, err := isPathExist(filepath)
	if err != nil {
		fmt.Println(err)
	} else {

		if !exist {

			file, err := os.Create(filepath)
			defer file.Close()
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	logFile, err := os.OpenFile(filepath, os.O_RDWR|os.O_APPEND, 0644)

	if err != nil {
		log.Fatalln("open file error !")
	}

	l.curLogFile = logFile
}

// 检查文件是否存在，不存在就新建
func (l *Log) checkLogFile(filepath string) {

	// 当前存在正写入的文件
	if l.curLogFile != nil {

		// 需要更换文件，加锁
		if !(strings.Compare(l.curLogFile.Name(), filepath) == 0) {

			l.mu.Lock()
			defer l.mu.Unlock()

			if strings.Compare(l.curLogFile.Name(), filepath) != 0 {
				l.curLogFile.Close()
				l.createNewFile(filepath)
			}
		}
	} else {
		l.createNewFile(filepath)
	}
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

// formatHeader
func (l *Log) formatHeader(buf *[]byte, t time.Time, prefix string) {

	*buf = append(*buf, prefix...)

	year, month, day := t.Date()
	itoa(buf, year, 4)
	*buf = append(*buf, '/')
	itoa(buf, int(month), 2)
	*buf = append(*buf, '/')
	itoa(buf, day, 2)
	*buf = append(*buf, ' ')

	hour, min, sec := t.Clock()
	itoa(buf, hour, 2)
	*buf = append(*buf, ':')
	itoa(buf, min, 2)
	*buf = append(*buf, ':')
	itoa(buf, sec, 2)

	*buf = append(*buf, '.')
	itoa(buf, t.Nanosecond()/1e6, 3)

	*buf = append(*buf, ' ')
}

// output writes the output for a logging event. The string s contains
func (l *Log) output(now time.Time, s string, prefix string) error {

	l.mu.Lock()
	defer l.mu.Unlock()
	l.buf = l.buf[:0]
	l.formatHeader(&l.buf, now, prefix)
	l.buf = append(l.buf, s...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	_, err := l.curLogFile.Write(l.buf)
	return err
}

// writeLog
func (l *Log) writeLog(msg string, level int) {

	t := time.Now()
	if l.logLevel <= level {

		l.checkLogFile(l.getLogPath(t))

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
		default:
			prefix = ""
		}

		l.output(t, msg, prefix)
	}
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
