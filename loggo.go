package loggo

import (
	"fmt"
	"strings"
)

const (
	ALL = iota
	DEBUG
	INFO
	ERROR
	FATAL
	OFF
)

type Logger struct {
	Level int `json:"level" ini:"level"` // log level

	// WithOutColor is the trigger of log color,if true ,the log will has no color.default is false(with color)
	WithOutColor bool        `json:"withOutColor" ini:"without_color"`
	writer       *FileWriter // file writer config
}

var (
	DefaultLog *Logger
)

func InitDefaultLog() {
	DefaultLog = &Logger{
		Level:        ALL,
		WithOutColor: false,
	}
	DefaultLog.SetWriter(NewDefaultWriter())
}
func (p *Logger) SetWriter(writer *FileWriter) {
	p.writer = writer
	p.writer.Init()
}

// 30（黑色）、31（红色）、32（绿色）、 33（黄色）、34（蓝色）、35（洋红）、36（青色）、37（白色）
func Text(color int, m ...string) {
	DefaultLog.txt(color, m...)
}

func Debugln(m ...string) {
	DefaultLog.logln(DEBUG, m...)
}
func Infoln(m ...string) {
	DefaultLog.logln(INFO, m...)
}
func Errorln(m ...string) {
	DefaultLog.logln(ERROR, m...)
}
func Fatalln(m ...string) {
	DefaultLog.logln(FATAL, m...)
	panic(strings.Join(m, ","))
}

func Debugfn(format string, args ...interface{}) {
	DefaultLog.logfn(DEBUG, format, args...)
}
func Infofn(format string, args ...interface{}) {
	DefaultLog.logfn(INFO, format, args...)
}
func Errorfn(format string, args ...interface{}) {
	DefaultLog.logfn(ERROR, format, args...)
}
func Fatalfn(format string, args ...interface{}) {
	DefaultLog.logfn(INFO, format, args...)
	panic(fmt.Sprintf(format, args...))
}

func PlainText(m ...string) {
	if len(m) == 0 {
		return
	}
	s := strings.Join(m, ",")
	DefaultLog.writer.Write([]byte(s))
}
func PlainTextln(m ...string) {
	if len(m) == 0 {
		return
	}
	s := strings.Join(m, ",")
	DefaultLog.writer.Write([]byte(s + "\n"))
}

func (p *Logger) Debugln(m ...string) {
	p.logln(DEBUG, m...)
}
func (p *Logger) Infoln(m ...string) {
	p.logln(INFO, m...)
}
func (p *Logger) Errorln(m ...string) {
	p.logln(ERROR, m...)
}
func (p *Logger) Fatalln(m ...string) {
	p.logln(FATAL, m...)
	panic(strings.Join(m, ","))
}

func (p *Logger) logln(level int, m ...string) {
	if p.Level > level {
		return
	}
	p.log(level, m...)
}

func (p *Logger) Debugfn(format string, args ...interface{}) {
	p.logfn(DEBUG, format, args...)
}
func (p *Logger) Infofn(format string, args ...interface{}) {
	p.logfn(INFO, format, args...)
}
func (p *Logger) Errorfn(format string, args ...interface{}) {
	p.logfn(ERROR, format, args...)
}
func (p *Logger) Fatalfn(format string, args ...interface{}) {
	p.logfn(INFO, format, args...)
	panic(fmt.Sprintf(format, args...))
}

func (p *Logger) logfn(level int, format string, args ...interface{}) {
	if p.Level > level {
		return
	}
	s := fmt.Sprintf(format, args...)
	p.log(level, []string{s}...)
}

// 30（黑色）、31（红色）、32（绿色）、 33（黄色）、34（蓝色）、35（洋红）、36（青色）、37（白色）
func getColor(level int) int {
	switch level {
	case DEBUG:
		return 33
	case INFO:
		return 32
	case ERROR:
		return 35
	case FATAL:
		return 31
	}
	return 30
}
func getLevelStr(level int) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO "
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	}
	return "UNKNOWN"
}
func (p *Logger) PlainText(m ...string) {
	if len(m) == 0 {
		return
	}
	s := strings.Join(m, ",")
	p.writer.Write([]byte(s))
}
func (p *Logger) PlainTextln(m ...string) {
	if len(m) == 0 {
		return
	}
	s := strings.Join(m, ",")
	p.writer.Write([]byte(s + "\n"))
}

// 30（黑色）、31（红色）、32（绿色）、 33（黄色）、34（蓝色）、35（洋红）、36（青色）、37（白色）
func (p *Logger) txt(color int, m ...string) {
	if len(m) == 0 {
		return
	}
	content := strings.Join(m, "\t")
	s := fmt.Sprintf("\033[%d;1m%s\033[0m", color, content)
	p.writer.Write([]byte(s))
}

// 30（黑色）、31（红色）、32（绿色）、 33（黄色）、34（蓝色）、35（洋红）、36（青色）、37（白色）
func (p *Logger) log(level int, m ...string) {
	if len(m) == 0 {
		return
	}
	content := strings.Join(m, "\t")
	var s string
	if p.WithOutColor {
		s = fmt.Sprintf("[%s|%s] %s\n", getTime(), getLevelStr(level), content)
	} else {
		s = fmt.Sprintf("\033[%d;1m[%s|%s] %s \033[0m\n", getColor(level), getTime(), getLevelStr(level), content)
	}
	p.writer.Write([]byte(s))
}

// Close implements io.Closer, and closes the current logfile.
func (l *Logger) Close() error {
	return l.writer.Close()
}
