package main

/*
 * @Description:
 * @Author: chenwei
 * @Date: 2020-01-15 16:47:52
 */

import (
	"github.com/silenceWe/loggo"
)

/*
 * @Description:
 * @Author: chenwei
 * @Date: 2020-01-15 16:47:52
 */

func main() {
	// default log
	loggo.InitDefaultLog(&loggo.LoggerOption{})
	loggo.Debugln("debug", "a", "b", "c")
	loggo.Infoln("info", "a")
	loggo.Errorln("error", "aaa", "b")

	// dynamic modify log level
	loggo.DefaultLogOption.Level = loggo.ERROR

	loggo.Debugfn("debug:%s,%d", "aaa", 123)
	loggo.Infofn("info:%s,%d", "aaa", 123)
	loggo.Errorfn("error:%s,%d", "aaa", 123)

	errLog := loggo.NewLoggo(&loggo.LoggerOption{
		RotateCron: "0 0 * * * *",   // 日志滚动定时任务cron表达式
		StdOut:     true,            // 是否在输出文件的同时，输出到标准输出
		Level:      loggo.ALL,       // 日志级别
		FileName:   "./log/err.log", // 日志文件名，滚动日志将在相同目录下生成
		MaxSize:    100,             // 单个日志最大size，超过阈值将会滚动

	})

	errLog.Infoln("err info:", "err")
	errLog.Infofn("info:%s,%d", "aaa", 123)

}
