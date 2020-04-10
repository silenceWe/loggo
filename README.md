## log for go

### default log
```
	loggo.InitDefaultLog(&loggo.LoggerOption{StdOut: true, Level: loggo.ALL})
	loggo.Debugln("debug")
	loggo.Infoln("info")
	loggo.Errorln("error")
    
```

### customer log

```
	log1 := loggo.NewLoggo(&loggo.LoggerOption{StdOut: true, Level: loggo.ALL, FileName: "./log/log1.log"})
	log1.Infoln("log1 info:", "log1")
```


### log config
```

        // RotateCron set the cron to rotate the log file
	RotateCron string

	// FileName is the file to write logs to.  Backup log files will be retained
	// in the same directory.  
	FileName string `json:"filename" ini:"filename"`

    // log level
	Level int

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

	// LocalTime determines if the time used for formatting the timestamps in
	// backup files is the computer's local time.  The default is to use UTC
	// time.
	LocalTime bool `json:"localtime" ini:"localtime"`

	// Compress determines if the rotated log files should be compressed
	// using gzip. The default is not to perform compression.
	Compress bool `json:"compress" ini:"compress"`

	// StdOut determines if the log should output to the std output
	StdOut bool `json:"stdOut" ini:"stdOut"`

	// WithOutColor is the trigger of log color,if true ,the log will has no color.default is false(with color)
	WithOutColor bool `json:"withOutColor" ini:"without_color"`
```