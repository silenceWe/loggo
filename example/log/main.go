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
	loggo.InitDefaultLog()
	loggo.Debugln("debug", "a", "b", "c")
	loggo.Infoln("info", "a")
	loggo.Errorln("error", "aaa", "b")

	// dynamic modify log level
	loggo.DefaultLog.Level = loggo.ERROR

	loggo.Debugfn("debug:%s,%d", "aaa", 123)
	loggo.Infofn("info:%s,%d", "aaa", 123)
	loggo.Errorfn("error:%s,%d", "aaa", 123)

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

}
