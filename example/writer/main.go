package main

import (
	"github.com/silenceWe/loggo"
	"github.com/silenceWe/loggo/writer"
	"log"
)

func main() {
	wtr1 := writer.NewLogWriter(&writer.WriterOption{StdOut: true, MaxAge: 7})
	wtr2 := writer.NewLogWriter(&writer.WriterOption{StdOut: true, MaxAge: 7, FileName: "./log2.log"})
	logger1 := log.Logger{}
	logger1.SetOutput(wtr1)
	logger1.SetFlags(log.LUTC)
	logger1.Println("test1")

	logger2 := log.Logger{}
	logger2.SetOutput(wtr2)
	logger2.SetFlags(log.LUTC)
	logger2.Println("test2")

	loggo1 := loggo.NewLoggo(&loggo.LoggerOption{
		RotateCron: "0 0 * * * *",   // 日志滚动定时任务cron表达式
		StdOut:     true,            // 是否在输出文件的同时，输出到标准输出
		Level:      loggo.ALL,       // 日志级别
		FileName:   "./log/err.log", // 日志文件名，滚动日志将在相同目录下生成
		MaxSize:    100,             // 单个日志最大size，超过阈值将会滚动

	})
}
