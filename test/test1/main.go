package main

import (
	"loggo/loggo"
)

func main() {
	loggo.InitDefaultLog(&loggo.LoggerOption{StdOut: true, Level: loggo.ALL})
	loggo.Debugln("debug")
	loggo.Infoln("info")
	loggo.Errorln("error")

	loggo.DefaultLogOption.Level = loggo.ERROR

	loggo.Debugfn("debug:%s,%d", "aaa", 123)
	loggo.Infofn("info:%s,%d", "aaa", 123)
	loggo.Errorfn("error:%s,%d", "aaa", 123)

	log1 := loggo.NewLoggo(&loggo.LoggerOption{StdOut: true, Level: loggo.ALL, FileName: "./log/log1.log"})
	log1.Infoln("log1 info:", "log1")

}
