package main

import (
	"github.com/silenceWe/loggo"
	"log"
	"time"
)

func main() {
	writer := &loggo.FileWriter{MaxAge: 7, FileName: "./log2.log", RotateCron: "0/5 * * * * *"}
	writer.Init()
	writer.SetBackupFormat("2006-01-02A150405000", "testformat-%s.log")
	logger1 := log.Logger{}
	logger1.SetOutput(writer)
	logger1.SetFlags(log.LUTC)
	logger1.Println("test1")

	loggo1 := &loggo.Logger{Level: loggo.ALL}
	loggo1.SetWriter(&loggo.FileWriter{FileName: "./loggo1.log"})
	loggo1.Infoln("test")

	for i := 0; i < 100; i++ {
		logger1.Println("test")
		time.Sleep(1 * time.Second)
	}
}
