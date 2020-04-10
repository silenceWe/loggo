package main

/*
 * @Description:
 * @Author: chenwei
 * @Date: 2020-01-15 16:47:52
 */

import "github.com/silenceWe/loggo"

/*
 * @Description:
 * @Author: chenwei
 * @Date: 2020-01-15 16:47:52
 */

func main() {
	loggo.InitDefaultLog(&loggo.LoggerOption{StdOut: true, Level: loggo.ALL, WithOutColor: true})
	loggo.Debugln("debug", "a", "b", "c")
	loggo.Infoln("info", "a")
	loggo.Errorln("error", "aaa", "b")

	loggo.DefaultLogOption.Level = loggo.ERROR

	loggo.Debugfn("debug:%s,%d", "aaa", 123)
	loggo.Infofn("info:%s,%d", "aaa", 123)
	loggo.Errorfn("error:%s,%d", "aaa", 123)

	log1 := loggo.NewLoggo(&loggo.LoggerOption{StdOut: true, Level: loggo.ALL, FileName: "./log/log1.log"})
	log1.Infoln("log1 info:", "log1")

}
