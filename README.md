## log for go
based on [lumberjack](https://github.com/natefinch/lumberjack)
### default log
```
	// default log
    	loggo.InitDefaultLog()
    	loggo.Debugln("debug", "a", "b", "c")
    	loggo.Infoln("info", "a")
    	loggo.Errorln("error", "aaa", "b")
    
    	// dynamic modify log level
    	loggo.DefaultLog.Level = loggo.ERROR
    
    	loggo.Debugfn("debug:%s,%d", "aaa", 123)
    	loggo.Infofn("info:%s,%d", "aaa", 123)
    	loggo.Errorfn("error:%s,%d", "aaa", 123)
    
```

### customer log

```
	
	writer := &loggo.FileWriter{MaxAge: 7, FileName: "./log2.log",RotateCron: "0/5 * * * * *"}
	writer.Init()
	writer.SetBackupFormat("2006-01-02A150405000","testformat-%s.log")
	logger1 := log.Logger{}
	logger1.SetOutput(writer)
	logger1.SetFlags(log.LUTC)
	logger1.Println("test1")

	errLog := &loggo.Logger{Level: loggo.ALL}
	errLog.SetWriter(&loggo.FileWriter{
		RotateCron: "0 0 * * * *",   // 日志滚动定时任务cron表达式
		FileName:   "./log/err.log", // 日志文件名，滚动日志将在相同目录下生成
		LocalTime:  false,
		MaxSize:    100,   // 单个日志最大size，超过阈值将会滚动
		MaxAge:     0,     // 日志保留时长(天)
		MaxBackups: 0,     // 日志保留个数
		Compress:   false, // 是否开启压缩
	})

	errLog.Infoln("err info:", "err")
	errLog.Infofn("info:%s,%d", "aaa", 123)
```

### use go built-in log package and loggo writer to rotate log file
```
	writer := &loggo.FileWriter{MaxAge: 7, FileName: "./log2.log"}
	writer.Init()

	logger1 := log.Logger{}
	logger1.SetOutput(writer)
	logger1.SetFlags(log.LUTC)
	logger1.Println("test1")
```


### log config
```

	Writer *FileWriter // file writer config
	Level  int         `json:"level" ini:"level"` // log level

	// WithOutColor is the trigger of log color,if true ,the log will has no color.default is false(with color)
	WithOutColor bool `json:"withOutColor" ini:"without_color"`
```
### writer config
```

	// RotateCron set the cron to rotate the log file
	RotateCron string `json:"rotate_cron" ini:"rotate_cron"`

	// FileName is the file to write logs to.  Backup log files will be retained
	// in the same directory.  It uses <processname>-lumberjack.log in
	// os.TempDir() if empty.
	FileName string `json:"filename" ini:"filename"`

	// LocalTime determines if the time used for formatting the timestamps in
	// backup files is the computer's local time.  The default is to use UTC
	// time.
	LocalTime bool `json:"localtime" ini:"localtime"`

	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes.
	MaxSize int `json:"maxsize" ini:"maxsize"`

	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	MaxAge int `json:"maxage" ini:"maxage"`

	// MaxBackups is the maximum number of old log files to retain.  The default
	// is to retain all old log files (though MaxAge may still cause them to get
	// deleted.)
	MaxBackups int `json:"maxbackups" ini:"maxbackups"`

	// Compress determines if the rotated log files should be compressed
	// using gzip. The default is not to perform compression.
	Compress bool `json:"compress" ini:"compress"`
```

