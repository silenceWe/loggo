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