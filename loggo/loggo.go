package loggo

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/robfig/cron"
)

const (
	printTimeFormat   = "2006-01-02 15:04:05.000"
	backupTimeFormat  = "2006-01-02T15-04-05.000"
	compressSuffix    = ".gz"
	defaultLogName    = "./log/default.log"
	defaultMaxSize    = 100
	defaultRotateCron = "0 0 0 * * *" // 00:00 AM every morning
)

// log level
const (
	ALL = iota
	DEBUG
	INFO
	ERROR
	FATAL
	OFF
)

// ensure we always implement io.WriteCloser
var _ io.WriteCloser = (*Logger)(nil)

type LoggerOption struct {
	// RotateCron set the cron to rotate the log file
	RotateCron string

	// FileName is the file to write logs to.  Backup log files will be retained
	// in the same directory.  It uses <processname>-lumberjack.log in
	// os.TempDir() if empty.
	FileName string `json:"filename" yaml:"filename"`

	Level int

	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes.
	MaxSize int `json:"maxsize" yaml:"maxsize"`

	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	MaxAge int `json:"maxage" yaml:"maxage"`

	// MaxBackups is the maximum number of old log files to retain.  The default
	// is to retain all old log files (though MaxAge may still cause them to get
	// deleted.)
	MaxBackups int `json:"maxbackups" yaml:"maxbackups"`

	// LocalTime determines if the time used for formatting the timestamps in
	// backup files is the computer's local time.  The default is to use UTC
	// time.
	LocalTime bool `json:"localtime" yaml:"localtime"`

	// Compress determines if the rotated log files should be compressed
	// using gzip. The default is not to perform compression.
	Compress bool `json:"compress" yaml:"compress"`

	// StdOut determines if the log should output to the std output
	StdOut bool `json:"stdOut" yaml:"stdOut"`
}
type Logger struct {
	option        *LoggerOption
	rotateRunning bool
	size          int64
	file          *os.File
	mu            sync.Mutex

	millCh    chan bool
	startMill sync.Once
}

var (
	// currentTime exists so it can be mocked out by tests.
	currentTime = time.Now

	// os_Stat exists so it can be mocked out by tests.
	os_Stat = os.Stat

	// megabyte is the conversion factor between MaxSize and bytes.  It is a
	// variable so tests can mock it out and not need to write megabytes of data
	// to disk.
	megabyte = 1024 * 1024

	defaultLog *Logger
)
var DefaultLogOption *LoggerOption

// init default log
func InitDefaultLog(option *LoggerOption) {
	DefaultLogOption = option
	if DefaultLogOption.FileName == "" {
		DefaultLogOption.FileName = defaultLogName
	}
	defaultLog = NewLoggo(DefaultLogOption)
}

func NewLoggo(option *LoggerOption) *Logger {
	if option.FileName == "" {
		panic("Please set the log file name")
	}
	option.LocalTime = true
	option.Compress = true
	if option.RotateCron == "" {
		option.RotateCron = defaultRotateCron
	}
	l := Logger{option: option}
	l.startRotateCron()
	return &l
}

func Debugln(m ...string) {
	defaultLog.println(DEBUG, m...)
}
func Infoln(m ...string) {
	defaultLog.println(INFO, m...)
}
func Errorln(m ...string) {
	defaultLog.println(ERROR, m...)
}
func Fatalln(m ...string) {
	defaultLog.println(FATAL, m...)
	panic(strings.Join(m, ","))
}

func Debugfn(format string, args ...interface{}) {
	defaultLog.printfn(DEBUG, format, args...)
}
func Infofn(format string, args ...interface{}) {
	defaultLog.printfn(INFO, format, args...)
}
func Errorfn(format string, args ...interface{}) {
	defaultLog.printfn(ERROR, format, args...)
}
func Fatalfn(format string, args ...interface{}) {
	defaultLog.printfn(INFO, format, args...)
	panic(fmt.Sprintf(format, args...))
}

func (p *Logger) Debugln(m ...string) {
	p.println(DEBUG, m...)
}
func (p *Logger) Infoln(m ...string) {
	p.println(INFO, m...)
}
func (p *Logger) Errorln(m ...string) {
	p.println(ERROR, m...)
}
func (p *Logger) Fatalln(m ...string) {
	p.println(FATAL, m...)
	panic(strings.Join(m, ","))
}

func (p *Logger) println(level int, m ...string) {
	if p.option.Level > level {
		return
	}
	p.printWithColor(level, m...)
}

func (p *Logger) Debugfn(format string, args ...interface{}) {
	p.printfn(DEBUG, format, args...)
}
func (p *Logger) Infofn(format string, args ...interface{}) {
	p.printfn(INFO, format, args...)
}
func (p *Logger) Errorfn(format string, args ...interface{}) {
	p.printfn(ERROR, format, args...)
}
func (p *Logger) Fatalfn(format string, args ...interface{}) {
	p.printfn(INFO, format, args...)
	panic(fmt.Sprintf(format, args...))
}

func (p *Logger) printfn(level int, format string, args ...interface{}) {
	if p.option.Level > level {
		return
	}
	s := fmt.Sprintf(format, args...)
	p.printWithColor(level, []string{s}...)
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
	if p.option.StdOut {
		log.Print(s)
	}
	p.Write([]byte(s))
}
func (p *Logger) PlainTextln(m ...string) {
	if len(m) == 0 {
		return
	}
	s := strings.Join(m, ",")
	if p.option.StdOut {
		fmt.Println(s)
	}
	p.Write([]byte(s + "\n"))
}

// 30（黑色）、31（红色）、32（绿色）、 33（黄色）、34（蓝色）、35（洋红）、36（青色）、37（白色）
func (p *Logger) printWithColor(level int, m ...string) {
	if len(m) == 0 {
		return
	}
	content := strings.Join(m, "\t")
	s := fmt.Sprintf("\033[%d;1m[%s|%s] %s \033[0m\n", getColor(level), getTime(), getLevelStr(level), content)
	if p.option.StdOut {
		os.Stdout.Write([]byte(s))
	}
	p.Write([]byte(s))
}
func getTime() string {
	return time.Now().Format(printTimeFormat)
}

// Write implements io.Writer.  If a write would cause the log file to be larger
// than MaxSize, the file is closed, renamed to include a timestamp of the
// current time, and a new log file is created using the original log file name.
// If the length of the write is greater than MaxSize, an error is returned.
func (l *Logger) Write(p []byte) (n int, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	writeLen := int64(len(p))
	if writeLen > l.max() {
		return 0, fmt.Errorf(
			"write length %d exceeds maximum file size %d", writeLen, l.max(),
		)
	}

	if l.file == nil {
		if err = l.openExistingOrNew(len(p)); err != nil {
			return 0, err
		}
	}

	if l.size+writeLen > l.max() {
		if err := l.rotate(); err != nil {
			return 0, err
		}
	}
	n, err = l.file.Write(p)
	l.size += int64(n)
	return n, err
}

// Close implements io.Closer, and closes the current logfile.
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.close()
}

// rotate
func (l *Logger) startRotateCron() {
	c := cron.New()
	l.Infoln("Log rotate cron:", l.option.RotateCron)
	c.AddFunc(l.option.RotateCron, func() {
		l.Infoln("------Start rotate log job")
		if l.rotateRunning {
			l.Infoln("job not finish wait...")
			return
		}
		l.rotateRunning = true
		if err := l.Rotate(); err != nil {
			l.Errorln("rotate error,", err.Error())
		}
		l.rotateRunning = false
	})
	c.Start()
}

// close closes the file if it is open.
func (l *Logger) close() error {
	if l.file == nil {
		return nil
	}
	err := l.file.Close()
	l.file = nil
	return err
}

// Rotate causes Logger to close the existing log file and immediately create a
// new one.  This is a helper function for applications that want to initiate
// rotations outside of the normal rotation rules, such as in response to
// SIGHUP.  After rotating, this initiates compression and removal of old log
// files according to the configuration.
func (l *Logger) Rotate() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.rotate()
}

// rotate closes the current file, moves it aside with a timestamp in the name,
// (if it exists), opens a new file with the original filename, and then runs
// post-rotation processing and removal.
func (l *Logger) rotate() error {
	if err := l.close(); err != nil {
		return err
	}
	if err := l.openNew(); err != nil {
		return err
	}
	l.mill()
	return nil
}

// openNew opens a new log file for writing, moving any old log file out of the
// way.  This methods assumes the file has already been closed.
func (l *Logger) openNew() error {
	err := os.MkdirAll(l.dir(), 0744)
	if err != nil {
		return fmt.Errorf("can't make directories for new logfile: %s", err)
	}

	name := l.filename()
	mode := os.FileMode(0644)
	info, err := os_Stat(name)
	if err == nil {
		// Copy the mode off the old logfile.
		mode = info.Mode()
		// move the existing file
		newname := backupName(name, l.option.LocalTime)
		if err := os.Rename(name, newname); err != nil {
			return fmt.Errorf("can't rename log file: %s", err)
		}

		// this is a no-op anywhere but linux
		if err := chown(name, info); err != nil {
			return err
		}
	}

	// we use truncate here because this should only get called when we've moved
	// the file ourselves. if someone else creates the file in the meantime,
	// just wipe out the contents.
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("can't open new logfile: %s", err)
	}
	l.file = f
	l.size = 0
	return nil
}

// backupName creates a new filename from the given name, inserting a timestamp
// between the filename and the extension, using the local time if requested
// (otherwise UTC).
func backupName(name string, local bool) string {
	dir := filepath.Dir(name)
	filename := filepath.Base(name)
	ext := filepath.Ext(filename)
	prefix := filename[:len(filename)-len(ext)]
	t := currentTime()
	if !local {
		t = t.UTC()
	}

	timestamp := t.Format(backupTimeFormat)
	return filepath.Join(dir, fmt.Sprintf("%s-%s%s", prefix, timestamp, ext))
}

// openExistingOrNew opens the logfile if it exists and if the current write
// would not put it over MaxSize.  If there is no such file or the write would
// put it over the MaxSize, a new file is created.
func (l *Logger) openExistingOrNew(writeLen int) error {
	l.mill()

	filename := l.filename()
	info, err := os_Stat(filename)
	if os.IsNotExist(err) {
		return l.openNew()
	}
	if err != nil {
		return fmt.Errorf("error getting log file info: %s", err)
	}

	if info.Size()+int64(writeLen) >= l.max() {
		return l.rotate()
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		// if we fail to open the old log file for some reason, just ignore
		// it and open a new log file.
		return l.openNew()
	}
	l.file = file
	l.size = info.Size()
	return nil
}

// genFilename generates the name of the logfile from the current time.
func (l *Logger) filename() string {
	if l.option.FileName != "" {
		return l.option.FileName
	}
	return "./log/default.log"
}

// millRunOnce performs compression and removal of stale log files.
// Log files are compressed if enabled via configuration and old log
// files are removed, keeping at most l.MaxBackups files, as long as
// none of them are older than MaxAge.
func (l *Logger) millRunOnce() error {
	if l.option.MaxBackups == 0 && l.option.MaxAge == 0 && !l.option.Compress {
		return nil
	}

	files, err := l.oldLogFiles()
	if err != nil {
		return err
	}

	var compress, remove []logInfo

	if l.option.MaxBackups > 0 && l.option.MaxBackups < len(files) {
		preserved := make(map[string]bool)
		var remaining []logInfo
		for _, f := range files {
			// Only count the uncompressed log file or the
			// compressed log file, not both.
			fn := f.Name()
			if strings.HasSuffix(fn, compressSuffix) {
				fn = fn[:len(fn)-len(compressSuffix)]
			}
			preserved[fn] = true

			if len(preserved) > l.option.MaxBackups {
				remove = append(remove, f)
			} else {
				remaining = append(remaining, f)
			}
		}
		files = remaining
	}
	if l.option.MaxAge > 0 {
		diff := time.Duration(int64(l.option.MaxAge) * int64(1*time.Minute))
		cutoff := currentTime().Add(-1 * diff)

		var remaining []logInfo
		for _, f := range files {
			if f.timestamp.Before(cutoff) {
				remove = append(remove, f)
			} else {
				remaining = append(remaining, f)
			}
		}
		files = remaining
	}

	if l.option.Compress {
		for _, f := range files {
			if !strings.HasSuffix(f.Name(), compressSuffix) {
				compress = append(compress, f)
			}
		}
	}

	for _, f := range remove {
		errRemove := os.Remove(filepath.Join(l.dir(), f.Name()))
		if err == nil && errRemove != nil {
			err = errRemove
		}
	}
	for _, f := range compress {
		fn := filepath.Join(l.dir(), f.Name())
		errCompress := compressLogFile(fn, fn+compressSuffix)
		if err == nil && errCompress != nil {
			err = errCompress
		}
	}

	return err
}

// millRun runs in a goroutine to manage post-rotation compression and removal
// of old log files.
func (l *Logger) millRun() {
	for _ = range l.millCh {
		// what am I going to do, log this?
		_ = l.millRunOnce()
	}
}

// mill performs post-rotation compression and removal of stale log files,
// starting the mill goroutine if necessary.
func (l *Logger) mill() {
	l.startMill.Do(func() {
		l.millCh = make(chan bool, 1)
		go l.millRun()
	})
	select {
	case l.millCh <- true:
	default:
	}
}

// oldLogFiles returns the list of backup log files stored in the same
// directory as the current log file, sorted by ModTime
func (l *Logger) oldLogFiles() ([]logInfo, error) {
	files, err := ioutil.ReadDir(l.dir())
	if err != nil {
		return nil, fmt.Errorf("can't read log file directory: %s", err)
	}
	logFiles := []logInfo{}

	prefix, ext := l.prefixAndExt()

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if t, err := l.timeFromName(f.Name(), prefix, ext); err == nil {
			logFiles = append(logFiles, logInfo{t, f})
			continue
		}
		if t, err := l.timeFromName(f.Name(), prefix, ext+compressSuffix); err == nil {
			logFiles = append(logFiles, logInfo{t, f})
			continue
		}
		// error parsing means that the suffix at the end was not generated
		// by lumberjack, and therefore it's not a backup file.
	}

	sort.Sort(byFormatTime(logFiles))

	return logFiles, nil
}

// timeFromName extracts the formatted time from the filename by stripping off
// the filename's prefix and extension. This prevents someone's filename from
// confusing time.parse.
func (l *Logger) timeFromName(filename, prefix, ext string) (time.Time, error) {
	if !strings.HasPrefix(filename, prefix) {
		return time.Time{}, errors.New("mismatched prefix")
	}
	if !strings.HasSuffix(filename, ext) {
		return time.Time{}, errors.New("mismatched extension")
	}
	ts := filename[len(prefix) : len(filename)-len(ext)]
	return time.Parse(backupTimeFormat, ts)
}

// max returns the maximum size in bytes of log files before rolling.
func (l *Logger) max() int64 {
	if l.option.MaxSize == 0 {
		return int64(defaultMaxSize * megabyte)
	}
	return int64(l.option.MaxSize) * int64(megabyte)
}

// dir returns the directory for the current filename.
func (l *Logger) dir() string {
	return filepath.Dir(l.filename())
}

// prefixAndExt returns the filename part and extension part from the Logger's
// filename.
func (l *Logger) prefixAndExt() (prefix, ext string) {
	filename := filepath.Base(l.filename())
	ext = filepath.Ext(filename)
	prefix = filename[:len(filename)-len(ext)] + "-"
	return prefix, ext
}

// compressLogFile compresses the given log file, removing the
// uncompressed log file if successful.
func compressLogFile(src, dst string) (err error) {
	f, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}
	defer f.Close()

	fi, err := os_Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat log file: %v", err)
	}

	if err := chown(dst, fi); err != nil {
		return fmt.Errorf("failed to chown compressed log file: %v", err)
	}

	// If this file already exists, we presume it was created by
	// a previous attempt to compress the log file.
	gzf, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, fi.Mode())
	if err != nil {
		return fmt.Errorf("failed to open compressed log file: %v", err)
	}
	defer gzf.Close()

	gz := gzip.NewWriter(gzf)

	defer func() {
		if err != nil {
			os.Remove(dst)
			err = fmt.Errorf("failed to compress log file: %v", err)
		}
	}()

	if _, err := io.Copy(gz, f); err != nil {
		return err
	}
	if err := gz.Close(); err != nil {
		return err
	}
	if err := gzf.Close(); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}
	if err := os.Remove(src); err != nil {
		return err
	}

	return nil
}

// logInfo is a convenience struct to return the filename and its embedded
// timestamp.
type logInfo struct {
	timestamp time.Time
	os.FileInfo
}

// byFormatTime sorts by newest time formatted in the name.
type byFormatTime []logInfo

func (b byFormatTime) Less(i, j int) bool {
	return b[i].timestamp.After(b[j].timestamp)
}

func (b byFormatTime) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b byFormatTime) Len() int {
	return len(b)
}
